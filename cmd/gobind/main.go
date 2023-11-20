package logger

// gomobile bind -target=android/arm,android/386 -javapkg rogo.iot.module .
import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	chclient "github.com/lthnh15032001/ngrok-impl/internal/client"
	"github.com/lthnh15032001/ngrok-impl/share/cos"
)

// triggerr
// R:5555:localhost:5555
func StreamLogToServer(server string, client string, findIpv4 bool) {
	config := chclient.Config{}
	config.Server = server

	if !findIpv4 {
		config.Remotes = []string{fmt.Sprintf("R:%s", client)}
	} else {
		ipv4, err := GetIpv4Address()
		if err != nil {
			log.Fatal("Can not find Ipv4 Address")
		}
		cPort, sPort, err := parsePorts(client)

		if err != nil {
			log.Fatal(err)
		}
		config.Remotes = []string{fmt.Sprintf("R:%d:%s:%d", cPort, ipv4, sPort)}
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
func parsePorts(input string) (cPort, sPort int, err error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid input format, expected 'cPort:sPort'")
		return
	}

	cPort, err = strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("failed to convert cPort to integer: %v", err)
		return
	}

	sPort, err = strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("failed to convert sPort to integer: %v", err)
		return
	}

	return
}
