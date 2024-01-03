package iot

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func newAuthCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Auth token",
		Long:  "Auth token",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			flags.Parse(args)
			args = flags.Args()
			if len(args) < 1 {
				log.Fatalf("A server and least one remote is required")
			}
			authTokenInput := args[0]
			h := os.Getenv("HOME")
			if h == "" {
				if u, err := user.Current(); err == nil {
					h = u.HomeDir
				}
			}
			c := filepath.Join(h, ".sndp-credentials", "sndp")
			// check folder exist then create
			if _, err := os.Stat(filepath.Join(h, ".sndp-credentials")); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir(filepath.Join(h, ".sndp-credentials"), os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
			}
			f, err := os.Create(c)
			if err != nil {
				log.Fatal(err)
			}
			// data := []byte(authTokenInput)
			if _, errWriteFile := f.WriteString(authTokenInput); errWriteFile != nil {
				log.Fatalf("Write credentials to file error %s", errWriteFile)
			}
			log.Printf("Write SNDP credentials at %s", c)
		},
	}
	return authCmd
}
