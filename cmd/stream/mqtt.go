package iot

import (
	"crypto/tls"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Broker         string
	Port           int
	ClientId       string
	ClientUsername string
	ClientPassword string
}

type Mqtt struct {
	*MqttConfig
	client mqtt.Client
}

func NewMqttConnection(c *MqttConfig) (*Mqtt, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtts://%s:%d", c.Broker, c.Port))
	tlsConfig := NewTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(c.ClientId)
	opts.SetUsername(c.ClientUsername)
	opts.SetPassword(c.ClientPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetMaxReconnectInterval(10 * time.Second)

	opts.SetAutoReconnect(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Println("...... mqtt reconnecting ......")
	})

	client := mqtt.NewClient(opts)

	mqttClient := &Mqtt{
		client: client,
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return mqttClient, nil
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected ")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (c *Mqtt) sub(topic string) {
	token := c.client.Subscribe(topic, 0, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func (c *Mqtt) publish() {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := c.client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func NewTlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}
