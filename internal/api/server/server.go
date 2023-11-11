package server

import (
	"log"

	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
)

func Init() {
	config := config.GetConfig()
	r := SetupRouter()
	err := r.Run(config.GetString("server.port"))
	if err != nil {
		log.Fatal(err)
	}
}
