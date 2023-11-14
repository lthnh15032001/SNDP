package iot

import (
	"sync"

	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/api/server"

	"github.com/spf13/cobra"
)

func newAPICommand() *cobra.Command {
	environment := "dev"
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Iot Streaming Log - Huda - ngrok streaming log ",
		Long:  "Iot Streaming Log - Huda - ngrok streaming log but long description",
		Run: func(cmd *cobra.Command, args []string) {
			config.Init(environment)
			var wg sync.WaitGroup
			errorCh := make(chan error)
			wg.Add(1)
			go func() {
				server.Init(errorCh)
				wg.Done()
			}()
			wg.Wait()
			close(errorCh)
			// return
		},
	}
	apiCmd.Flags().StringVar(&environment, "e", "", "")
	return apiCmd
}
