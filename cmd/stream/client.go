package iot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	chclient "github.com/lthnh15032001/ngrok-impl/internal/client"
	"github.com/lthnh15032001/ngrok-impl/share/cos"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newClientTCPCommand() *cobra.Command {

	config := &chclient.Config{Headers: http.Header{}}
	clientCmd := &cobra.Command{
		Use:           "client",
		Short:         "Iot Streaming TCP Log - Huda",
		Long:          "Iot Streaming TCP Log - Huda",
		SilenceErrors: true,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client(args, cmd.Flags(), config)
		},
	}
	clientCmd.Flags().StringVar(&config.Fingerprint, "fingerprint", "", "")
	clientCmd.Flags().StringVar(&config.Auth, "auth", "", "")
	clientCmd.Flags().DurationVar(&config.KeepAlive, "keepalive", 25*time.Second, "")
	clientCmd.Flags().IntVar(&config.MaxRetryCount, "max-retry-count", -1, "")
	clientCmd.Flags().DurationVar(&config.MaxRetryInterval, "max-retry-interval", 0, "")
	clientCmd.Flags().StringVar(&config.Proxy, "proxy", "", "")
	clientCmd.Flags().StringVar(&config.TLS.CA, "tls-ca", "", "")
	clientCmd.Flags().BoolVar(&config.TLS.SkipVerify, "tls-skip-verify", false, "")
	clientCmd.Flags().StringVar(&config.TLS.Cert, "tls-cert", "", "")
	clientCmd.Flags().StringVar(&config.TLS.Key, "tls-key", "", "")
	clientCmd.Flags().Var(&headerFlags{config.Headers}, "header", "")

	clientCmd.Flags().String("hostname", "", "")
	clientCmd.Flags().String("sni", "", "")
	clientCmd.Flags().Bool("pid", false, "")
	clientCmd.Flags().Bool("v", false, "")
	return clientCmd
}

type headerFlags struct {
	http.Header
}

// Type implements pflag.Value.
func (*headerFlags) Type() string {
	panic("unimplemented")
}

func (flag *headerFlags) String() string {
	out := ""
	for k, v := range flag.Header {
		out += fmt.Sprintf("%s: %s\n", k, v)
	}
	return out
}

func (flag *headerFlags) Set(arg string) error {
	index := strings.Index(arg, ":")
	if index < 0 {
		return fmt.Errorf(`Invalid header (%s). Should be in the format "HeaderName: HeaderContent"`, arg)
	}
	if flag.Header == nil {
		flag.Header = http.Header{}
	}
	key := arg[0:index]
	value := arg[index+1:]
	flag.Header.Set(key, strings.TrimSpace(value))
	return nil
}

func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	if err := os.WriteFile("github.com/lthnh15032001/ngrok-impl.pid", pid, 0644); err != nil {
		log.Fatal(err)
	}
}
func client(args []string, flags *pflag.FlagSet, config *chclient.Config) {

	hostname := flags.Lookup("hostname").Value.String()
	sni := flags.Lookup("sni").Value.String()
	pid := strToBoolean(flags.Lookup("pid").Value.String())
	// verbose := strToBoolean(flags.Lookup("verbose").Value.String())

	flags.Parse(args)
	//pull out options, put back remaining args
	args = flags.Args()
	// if len(args) < 1 {
	// 	log.Fatalf("At Least one remote is required")
	// }
	// config.Server = "https://tunnel.rogo.com.vn"
	// config.Remotes = args[0:]
	if len(args) < 2 {
		log.Fatalf("A server and least one remote is required")
	}
	config.Server = args[0]
	config.Remotes = args[1:]
	//default auth
	if config.Auth == "" {
		config.Auth = os.Getenv("AUTH")
	}
	//move hostname onto headers
	if hostname != "" {
		config.Headers.Set("Host", hostname)
		config.TLS.ServerName = hostname
	}

	if sni != "" {
		config.TLS.ServerName = sni
	}
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get(%s) failed\n", url)
	}
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll failed")
	}
	config.Headers.Set("X-FromIPAddress", fmt.Sprintf("%s", ip))
	config.Headers.Set("X-Runtime", runtime.GOOS)

	//ready
	c, err := chclient.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}
	// TODO: if there is no err, then saved that session to
	c.Debug = true
	if pid {
		generatePidFile()
	}
	go cos.GoStats()
	ctx := cos.InterruptContext()
	if err := c.Start(ctx); err != nil {
		log.Fatal(err)
	}
	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}
}
