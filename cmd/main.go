package main

import (
	iot "iot-stream/cmd/stream"
	"iot-stream/internal/logging"
	"iot-stream/internal/signals"
	"os"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := iot.Execute(ctx); err != nil {
		logging.LoggerFromContext(ctx).Error(err)
		os.Exit(1)
	}
}
