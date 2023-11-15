package iot

import (
	"fmt"
	"log"
	"net"
	"os"

	chclient "github.com/lthnh15032001/ngrok-impl/internal/client"
	"github.com/lthnh15032001/ngrok-impl/share/cos"
)

func StreamLogToServer(server string, port int, findIpv4 bool, options ...string) {
	config := chclient.Config{}
	config.Server = server

	if !findIpv4 {
		var remote string
		if len(options) < 1 {
			log.Fatal("Missing your ipv4, use findIpv4 is true instead")
		}
		for _, opt := range options {
			remote = opt
			config.Remotes = []string{fmt.Sprintf("R:%s:%d", remote, port)}
		}

	} else {
		ipv4, err := GetIpv4Address()
		if err != nil {
			log.Fatal("Can not find Ipv4 Address")
		}
		config.Remotes = []string{fmt.Sprintf("R:%s:%d", ipv4, port)}
	}
	//default auth
	if config.Auth == "" {
		config.Auth = os.Getenv("AUTH")
	}

	//ready
	c, err := chclient.NewClient(&config)

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

func GetIpv4Address() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				fmt.Printf("Your IPv4 Address: %s \n", ipNet.IP.String())
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("IPv4 address not found")
}
