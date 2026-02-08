package agent

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Miklakapi/gometrum/internal/mqtt"
	"github.com/Miklakapi/gometrum/internal/sensors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type agent struct {
	client *mqtt.MQTTClient

	groupedSensors map[time.Duration][]sensors.Sensor

	stateBase         string
	discoveryBase     string
	availabilityTopic string

	deviceId     string
	deviceName   string
	manufacturer string
	model        string
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

func New(s Settings, sens []sensors.Sensor) (*agent, error) {
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

	stateBase := s.StatePrefix + "/" + s.DeviceId
	availabilityTopic := stateBase + "/availability"

	o.SetWill(availabilityTopic, "offline", 1, true)

	client := mqtt.New(o)
	client.SetAvailability(availabilityTopic, []byte("online"))

	return &agent{
		client: client,

		groupedSensors: groupByInterval(sens),

		stateBase:         stateBase,
		discoveryBase:     s.DiscoveryPrefix,
		availabilityTopic: availabilityTopic,

		deviceId:     s.DeviceId,
		deviceName:   s.DeviceName,
		manufacturer: s.Manufacturer,
		model:        s.Model,
	}, nil
}

func (a *agent) Run(ctx context.Context) error {
	if err := a.client.Connect(10 * time.Second); err != nil {
		return err
	}
	defer func() {
		if a.availabilityTopic != "" {
			if err := a.client.Publish(a.availabilityTopic, 1, true, []byte("offline")); err != nil {
				slog.Warn("mqtt publish offline failed", "topic", a.availabilityTopic, "err", err)
			}
		}
		a.client.Close()
	}()

	if err := a.publishDiscovery(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for interval, group := range a.groupedSensors {
		wg.Go(func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			sensorsStateCache := make(map[string]string, len(group))

			a.collectAndPublishGroup(ctx, group, sensorsStateCache)
			for {
				select {
				case <-ticker.C:
					a.collectAndPublishGroup(ctx, group, sensorsStateCache)
				case <-ctx.Done():
					return
				}
			}
		})
	}

	wg.Wait()
	return nil
}

func (a *agent) Purge() error {
	if err := a.client.Connect(10 * time.Second); err != nil {
		return err
	}
	defer a.client.Close()

	for _, group := range a.groupedSensors {
		for _, s := range group {
			topic := fmt.Sprintf("%s/sensor/%s/%s/config", a.discoveryBase, a.deviceId, s.Key())
			if err := a.client.Publish(topic, 1, true, []byte{}); err != nil {
				return fmt.Errorf("purge: clear discovery failed (topic=%s): %w", topic, err)
			}
		}
	}

	for _, group := range a.groupedSensors {
		for _, s := range group {
			topic := fmt.Sprintf("%s/%s/state", a.stateBase, s.Key())
			if err := a.client.Publish(topic, 1, true, []byte{}); err != nil {
				return fmt.Errorf("purge: clear state failed (topic=%s): %w", topic, err)
			}
		}
	}

	if a.availabilityTopic != "" {
		if err := a.client.Publish(a.availabilityTopic, 1, true, []byte{}); err != nil {
			slog.Warn("purge: clear availability failed", "topic", a.availabilityTopic, "err", err)
		}
	}

	return nil
}

func groupByInterval(list []sensors.Sensor) map[time.Duration][]sensors.Sensor {
	groups := make(map[time.Duration][]sensors.Sensor)

	for _, s := range list {
		groups[s.Interval()] = append(groups[s.Interval()], s)
	}

	return groups
}
