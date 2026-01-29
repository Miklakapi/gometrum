package agent

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/Miklakapi/gometrum/internal/mqtt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type agent struct {
	client *mqtt.MQTTClient
}

type Settings struct {
	Host            string
	Port            int
	Username        string
	Password        string
	ClientID        string
	DiscoveryPrefix string
	StatePrefix     string

	DeviceId     string
	DeviceName   string
	Manufacturer string
	Model        string
}

func New(s Settings) (*agent, error) {
	o := MQTT.NewClientOptions()
	addr := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
	o.AddBroker("tcp://" + addr)

	o.SetClientID(s.ClientID)
	o.SetUsername(s.Username)
	o.SetPassword(s.Password)

	o.SetCleanSession(true)
	o.SetAutoReconnect(true)
	o.SetConnectRetry(true)

	o.SetConnectTimeout(10 * time.Second)
	o.SetKeepAlive(30 * time.Second)
	o.SetPingTimeout(10 * time.Second)

	client := mqtt.New(o)

	return &agent{
		client: client,
	}, nil
}

func (a *agent) Run(ctx context.Context) error {
	if err := a.client.Connect(10 * time.Second); err != nil {
		return err
	}
	defer a.client.Close()

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			slog.Info("Test")
			// a.mqtt.Publish(...)
		}
	}
}
