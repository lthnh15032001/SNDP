package main

import (
	"iot-stream/internal/logging"
	"iot-stream/internal/signals"
	"os"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := Execute(ctx); err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
}
