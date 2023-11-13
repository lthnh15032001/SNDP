package main

import (
	"os"

	iot "github.com/lthnh15032001/ngrok-impl/cmd/stream"
	"github.com/lthnh15032001/ngrok-impl/internal/logging"
	"github.com/lthnh15032001/ngrok-impl/internal/signals"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := iot.Execute(ctx); err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
}
