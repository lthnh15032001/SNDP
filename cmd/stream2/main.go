package iot

import (
	chclient "iot-stream/internal/client"
	"iot-stream/share/cos"
	"log"
	"os"
)
func Client(config *chclient.Config) {
	config.Server = "https://tunnel.rogo.com.vn"
	config.Remotes = []string{"R:localhost:5555"}
	//default auth
	if config.Auth == "" {
		config.Auth = os.Getenv("AUTH")
	}

	//ready
	c, err := chclient.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}
	c.Debug = true
	go cos.GoStats()
	ctx := cos.InterruptContext()
	if err := c.Start(ctx); err != nil {
		log.Fatal(err)
	}
	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}
}
