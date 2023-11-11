package main

import (
	"os"

	iot "github.com/lthnh15032001/ngrok-impl/cmd/stream"
	iot2 "github.com/lthnh15032001/ngrok-impl/cmd/stream2"
	chclient "github.com/lthnh15032001/ngrok-impl/internal/client"
	"github.com/lthnh15032001/ngrok-impl/internal/logging"
	"github.com/lthnh15032001/ngrok-impl/internal/signals"
)

func main() {
	ctx := signals.SetupSignalHandler()
	iot2.Client(&chclient.Config{})
	if err := iot.Execute(ctx); err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
}
