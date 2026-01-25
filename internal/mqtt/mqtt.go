package mqtt

import (
	"log/slog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client MQTT.Client
}

func New(o *MQTT.ClientOptions) (*MQTTClient, error) {
	o.OnConnect = func(c MQTT.Client) {
		slog.Info("mqtt connected")
	}

	o.OnConnectionLost = func(c MQTT.Client, err error) {
		slog.Error("mqtt connection lost", "err", err)
	}

	o.OnReconnecting = func(c MQTT.Client, co *MQTT.ClientOptions) {
		slog.Info("mqtt reconnecting to", "servers", co.Servers)
	}

	client := MQTT.NewClient(o)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTClient{
		client: client,
	}, nil
}

func (m *MQTTClient) Publish(topic string, qos byte, retain bool, payload []byte) error {
	token := m.client.Publish(topic, qos, retain, payload)
	token.Wait()
	return token.Error()
}

func (m *MQTTClient) Close() {
	m.client.Disconnect(250)
}
