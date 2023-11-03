package server

import (
	"iot-stream/internal/api/config"
	"log"
)

func Init() {
	config := config.GetConfig()
	r := SetupRouter()
	err := r.Run(config.GetString("server.port"))
	if err != nil {
		log.Fatal(err)
	}
}
