package chserver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lthnh15032001/ngrok-impl/internal/store"
	chshare "github.com/lthnh15032001/ngrok-impl/share"
	"github.com/lthnh15032001/ngrok-impl/share/ccrypto"
	"github.com/lthnh15032001/ngrok-impl/share/cio"
	"github.com/lthnh15032001/ngrok-impl/share/cnet"
	"github.com/lthnh15032001/ngrok-impl/share/settings"

	"github.com/gorilla/websocket"
	"github.com/jpillora/requestlog"
	"golang.org/x/crypto/ssh"
)

// Config is the configuration for the chisel service
type Config struct {
	KeySeed    string
	KeyFile    string
	AuthFile   string
	Auth       string
	Proxy      string
	Socks5     bool
	Reverse    bool
	KeepAlive  time.Duration
	TLS        TLSConfig
	WithDBAuth bool
}

// Server respresent a chisel service
type Server struct {
	*cio.Logger
	config       *Config
	fingerprint  string
	httpServer   *cnet.HTTPServer
	reverseProxy *httputil.ReverseProxy
	sessCount    int32
	sessions     *settings.Users
	sshConfig    *ssh.ServerConfig
	users        *settings.UserIndex
	sc           store.Interface
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  settings.EnvInt("WS_BUFF_SIZE", 0),
	WriteBufferSize: settings.EnvInt("WS_BUFF_SIZE", 0),
}

// NewServer creates and returns a new chisel server
func NewServer(c *Config) (*Server, error) {
	server := &Server{
		config:     c,
		httpServer: cnet.NewHTTPServer(),
		Logger:     cio.NewLogger("server"),
		sessions:   settings.NewUsers(),
	}
	server.Info = true
	server.Debug = true
	server.users = settings.NewUserIndex(server.Logger)
	// if db auth user acl then lookup user to authen user, else load user from file, default is non authen
	if !c.WithDBAuth {
		// auth user from a file => auth user ACL from database
		if c.AuthFile != "" {
			if err := server.users.LoadUsers(c.AuthFile); err != nil {
				return nil, err
			}
		}
		if c.Auth != "" {
			u := &settings.User{Addrs: []*regexp.Regexp{settings.UserAllowAll}}
			u.Name, u.Pass = settings.ParseAuth(c.Auth)
			if u.Name != "" {
				server.users.AddUser(u)
			}
		}
	}
	var pemBytes []byte
	var err error
	if c.KeyFile != "" {
		var key []byte

		if ccrypto.IsChiselKey([]byte(c.KeyFile)) {
			key = []byte(c.KeyFile)
		} else {
			key, err = os.ReadFile(c.KeyFile)
			if err != nil {
				log.Fatalf("Failed to read key file %s", c.KeyFile)
			}
		}

		pemBytes = key
		if ccrypto.IsChiselKey(key) {
			pemBytes, err = ccrypto.ChiselKey2PEM(key)
			if err != nil {
				log.Fatalf("Invalid key %s", string(key))
			}
		}
	} else {
		//generate private key (optionally using seed)
		pemBytes, err = ccrypto.Seed2PEM(c.KeySeed)
		if err != nil {
			log.Fatal("Failed to generate key")
		}
	}

	//convert into ssh.PrivateKey
	private, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatal("Failed to parse key")
	}
	//fingerprint this key
	server.fingerprint = ccrypto.FingerprintKey(private.PublicKey())
	//create ssh config
	server.sshConfig = &ssh.ServerConfig{
		ServerVersion:    "SSH-" + chshare.ProtocolVersion + "-server",
		PasswordCallback: server.authUser,
	}
	server.sshConfig.AddHostKey(private)
	//setup reverse proxy
	if c.Proxy != "" {
		u, err := url.Parse(c.Proxy)
		if err != nil {
			return nil, err
		}
		if u.Host == "" {
			return nil, server.Errorf("Missing protocol (%s)", u)
		}
		server.reverseProxy = httputil.NewSingleHostReverseProxy(u)
		//always use proxy host
		server.reverseProxy.Director = func(r *http.Request) {
			//enforce origin, keep path
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.Host = u.Host
		}
	}
	//print when reverse tunnelling is enabled
	if c.Reverse {
		server.Infof("Reverse tunnelling enabled")
	}
	return server, nil
}

// Run is responsible for starting the chisel service.
// Internally this calls Start then Wait.
func (s *Server) Run(host, port string) error {
	if err := s.Start(host, port); err != nil {
		return err
	}
	return s.Wait()
}

// Start is responsible for kicking off the http server
func (s *Server) Start(host, port string) error {
	return s.StartContext(context.Background(), host, port)
}

// StartContext is responsible for kicking off the http server & mysql server,
// and can be closed by cancelling the provided context
func (s *Server) StartContext(ctx context.Context, host, port string) error {
	s.Infof("Fingerprint %s", s.fingerprint)
	var err error
	if s.config.WithDBAuth {
		s.sc, _, err = store.GetOnce()
	}
	if err != nil {
		return err
	}
	if s.users.Len() > 0 {
		s.Infof("User authentication enabled")
	}
	if s.reverseProxy != nil {
		s.Infof("Reverse proxy enabled")
	}
	l, err := s.listener(host, port)
	if err != nil {
		return err
	}

	// broadcast adb log
	// go s.broadcastLogcatOutput()
	// http.HandleFunc("/ws", s.handleWebSocketConnection)
	h := http.Handler(http.HandlerFunc(s.handleClientHandler))

	if s.Debug {
		o := requestlog.DefaultOptions
		o.TrustProxy = true
		h = requestlog.WrapWith(h, o)
	}
	return s.httpServer.GoServe(ctx, l, h)
}

// Wait waits for the http server to close
func (s *Server) Wait() error {
	return s.httpServer.Wait()
}

// Close forcibly closes the http server
func (s *Server) Close() error {
	return s.httpServer.Close()
}

// GetFingerprint is used to access the server fingerprint
func (s *Server) GetFingerprint() string {
	return s.fingerprint
}

func (s *Server) GetUsernameAndUserId(user string) (string, string) {
	res1 := strings.Split(user, ".")
	// if len(res1) != 2 {
	// 	errors.New("invalid authentication for user: %s")
	// }
	return res1[0], res1[1]
}

// authUser is responsible for validating the ssh user / password combination
func (s *Server) authUser(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	// fmt.Printf("ccccccccccccccc %s", c.User())
	username, userId := s.GetUsernameAndUserId(c.User())
	if s.config.WithDBAuth {
		sc := s.sc
		//TODO: do authen for developer to get userId
		checkUser, err := sc.CheckUserExist(username, userId)
		if err != nil {
			return nil, errors.New("invalid authentication for username: %s")
		}
		if checkUser.Password != string(password) {
			s.Debugf("Login failed for userss: %s", username)
			return nil, errors.New("invalid authentication for username: %s")
		}
		var policies []string
		err = json.Unmarshal(checkUser.UserRemotePolicy, &policies)
		if err != nil {
			return nil, errors.New("invalid authentication because of policies in wrong format for username: %s")
		}
		s.AddUser(checkUser.Username, checkUser.Password, policies...)
		user, _ := s.users.Get(username)

		s.sessions.Set(string(c.SessionID()), user)
		return nil, nil
	}
	// check if user authentication is enabled and if not, allow all
	if s.users.Len() == 0 {
		return nil, nil
	}
	// check the user exists and has matching password
	n := username
	user, found := s.users.Get(n)
	if !found || user.Pass != string(password) {
		s.Debugf("Login failed for user: %s", n)
		return nil, errors.New("invalid authentication for username: %s")
	}
	// insert the user session map
	// TODO this should probably have a lock on it given the map isn't thread-safe
	s.sessions.Set(string(c.SessionID()), user)
	return nil, nil
}

// AddUser adds a new user into the server user index
func (s *Server) AddUser(user, pass string, addrs ...string) error {
	authorizedAddrs := []*regexp.Regexp{}
	for _, addr := range addrs {
		s.Infof("addr", addr)
		authorizedAddr, err := regexp.Compile(addr)
		if err != nil {
			return err
		}
		authorizedAddrs = append(authorizedAddrs, authorizedAddr)
	}
	s.users.AddUser(&settings.User{
		Name:  user,
		Pass:  pass,
		Addrs: authorizedAddrs,
	})
	return nil
}

// DeleteUser removes a user from the server user index
func (s *Server) DeleteUser(user string) {
	s.users.Del(user)
}

// ResetUsers in the server user index.
// Use nil to remove all.
func (s *Server) ResetUsers(users []*settings.User) {
	s.users.Reset(users)
}
