package mqtt

import (
	"fmt"
	"log/slog"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client MQTT.Client
}

func New(o *MQTT.ClientOptions) *MQTTClient {
	o.OnConnect = func(c MQTT.Client) {
		slog.Info("mqtt connected")
	}

	o.OnConnectionLost = func(c MQTT.Client, err error) {
		slog.Warn("mqtt connection lost", "err", err)
	}

	o.OnReconnecting = func(c MQTT.Client, co *MQTT.ClientOptions) {
		slog.Info("mqtt reconnecting to", "servers", co.Servers)
	}

	client := MQTT.NewClient(o)

	return &MQTTClient{
		client: client,
	}
}

func (m *MQTTClient) Connect(timeout time.Duration) error {
	token := m.client.Connect()
	if !token.WaitTimeout(timeout) {
		return fmt.Errorf("mqtt connect timeout after %s", timeout)
	}
	return token.Error()
}

func (m *MQTTClient) Publish(topic string, qos byte, retain bool, payload []byte) error {
	token := m.client.Publish(topic, qos, retain, payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("mqtt publish timeout (topic=%s)", topic)
	}
	return token.Error()
}

func (m *MQTTClient) Close() {
	m.client.Disconnect(250)
}
