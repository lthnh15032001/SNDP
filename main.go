package main

import (
	"flag"
	"fmt"
	"iot-stream/internal/api/config"
	"iot-stream/internal/api/server"
	"os"
)

func main() {
	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	server.Init()
}
