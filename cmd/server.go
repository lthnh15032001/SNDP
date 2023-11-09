package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	chserver "iot-stream/internal/server"
	"iot-stream/share/ccrypto"
	"iot-stream/share/cos"
	"iot-stream/share/settings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newServerTCPCommand() *cobra.Command {
	config := &chserver.Config{}
	clientCmd := &cobra.Command{
		Use:           "server",
		Short:         "Iot Streaming TCP Log - Huda",
		Long:          "Iot Streaming TCP Log - Huda",
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			serverTCP(args, cmd.Flags(), config)
		},
	}
	clientCmd.Flags().StringVar(&config.KeySeed, "key", "", "")
	clientCmd.Flags().StringVar(&config.KeyFile, "keyfile", "", "")
	clientCmd.Flags().StringVar(&config.AuthFile, "authfile", "", "")
	clientCmd.Flags().StringVar(&config.Auth, "auth", "", "")
	clientCmd.Flags().DurationVar(&config.KeepAlive, "keepalive", 25*time.Second, "")
	clientCmd.Flags().StringVar(&config.Proxy, "proxy", "", "")
	clientCmd.Flags().StringVar(&config.Proxy, "backend", "", "")
	clientCmd.Flags().BoolVar(&config.Socks5, "socks5", false, "")
	clientCmd.Flags().BoolVar(&config.Reverse, "reverse", false, "")
	clientCmd.Flags().StringVar(&config.TLS.Key, "tls-key", "", "")
	clientCmd.Flags().StringVar(&config.TLS.Cert, "tls-cert", "", "")
	clientCmd.Flags().Var(multiFlag{&config.TLS.Domains}, "tls-domain", "")
	clientCmd.Flags().StringVar(&config.TLS.CA, "tls-ca", "", "")

	clientCmd.Flags().String("host", os.Getenv("HOST"), "")
	clientCmd.Flags().String("p", "0.0.0.0", "")
	clientCmd.Flags().String("port", os.Getenv("PORT"), "")
	clientCmd.Flags().Bool("pid", false, "")
	clientCmd.Flags().Bool("v", false, "")
	clientCmd.Flags().String("keygen", "", "")
	return clientCmd
}

type multiFlag struct {
	values *[]string
}

// Type implements pflag.Value.
func (multiFlag) Type() string {
	panic("unimplemented")
}

func (flag multiFlag) String() string {
	return strings.Join(*flag.values, ", ")
}

func (flag multiFlag) Set(arg string) error {
	*flag.values = append(*flag.values, arg)
	return nil
}

func strToBoolean(str string) bool {
	boo, err := strconv.ParseBool(str)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error:", err)
	}
	return boo
}
func serverTCP(args []string, flags *pflag.FlagSet, config *chserver.Config) {

	host := flags.Lookup("host").Value.String()
	p := flags.Lookup("p").Value.String()
	port := flags.Lookup("port").Value.String()
	pid := strToBoolean(flags.Lookup("pid").Value.String())
	keyGen := flags.Lookup("keygen").Value.String()

	flags.Parse(args)

	if keyGen != "" {
		if err := ccrypto.GenerateKeyFile(keyGen, config.KeySeed); err != nil {
			log.Fatal(err)
		}
		return
	}

	if config.KeySeed != "" {
		log.Print("Option `--key` is deprecated and will be removed in a future version of chisel.")
		log.Print("Please use `chisel server --keygen /file/path`, followed by `chisel server --keyfile /file/path` to specify the SSH private key")
	}

	if host == "" {
		host = os.Getenv("HOST")
	}
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = p
	}
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}
	if config.KeyFile == "" {
		config.KeyFile = settings.Env("KEY_FILE")
	} else if config.KeySeed == "" {
		config.KeySeed = settings.Env("KEY")
	}
	s, err := chserver.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	s.Debug = true
	if pid {
		generatePidFile()
	}
	go cos.GoStats()
	ctx := cos.InterruptContext()
	if err := s.StartContext(ctx, host, port); err != nil {
		log.Fatal(err)
	}
	if err := s.Wait(); err != nil {
		log.Fatal(err)
	}
}
