package iot

import (
	"context"
	"encoding/json"
	"fmt"

	versionpkg "github.com/lthnh15032001/ngrok-impl/internal/version"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:               "huda",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Run: func(cmd *cobra.Command, args []string) {
			version := versionpkg.GetVersion()
			log.WithFields(log.Fields{
				"version": version.Version,
				"commit":  version.GitCommit,
			}).Info("Starting github.com/lthnh15032001/ngrok-impl")
			cmd.HelpFunc()(cmd, args)
		},
	}
)

func Execute(ctx context.Context) error {
	rootCmd.AddCommand(newAPICommand()) // new API server

	rootCmd.AddCommand(newHttpCommand())      // http streaming
	rootCmd.AddCommand(newClientTCPCommand()) // tcp streaming
	rootCmd.AddCommand(newServerTCPCommand()) // tcp streaming

	// define version
	rootCmd.AddCommand(newVersionCommand())
	return rootCmd.ExecuteContext(ctx)
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := json.Marshal(versionpkg.GetVersion())
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		},
	}
}
