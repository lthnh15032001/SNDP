package iot

import (
	"flag"
	"fmt"
	"os"

	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/api/server"

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
