package main

import (
	"flag"
	"fmt"
	"iot-stream/internal/api/config"
	"iot-stream/internal/api/server"
	"os"

	"github.com/spf13/cobra"
)

func newAPICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Iot Streaming Log - Huda - ngrok streaming log ",
		Long:  "Iot Streaming Log - Huda - ngrok streaming log but long description",
		Run: func(cmd *cobra.Command, args []string) {
			environment := flag.String("e", "dev", "")
			flag.Usage = func() {
				fmt.Println("Usage: server -e {mode}")
				os.Exit(1)
			}
			flag.Parse()
			config.Init(*environment)
			server.Init()
		},
	}

}
