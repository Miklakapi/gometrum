package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Miklakapi/gometrum/internal/mqtt"
	"github.com/Miklakapi/gometrum/internal/sensors"
)

type agent struct {
	pub mqtt.Publisher

	groupedSensors map[time.Duration][]sensors.Sensor

	stateBase         string
	discoveryBase     string
	availabilityTopic string

	deviceId     string
	deviceName   string
	manufacturer string
	model        string

	once bool
}

type Settings struct {
	DiscoveryPrefix string
	StatePrefix     string

	DeviceId     string
	DeviceName   string
	Manufacturer string
	Model        string

	Once bool
}

func New(s Settings, sens []sensors.Sensor, pub mqtt.Publisher) (*agent, error) {
	stateBase := s.StatePrefix + "/" + s.DeviceId
	availabilityTopic := stateBase + "/availability"

	pub.SetAvailability(availabilityTopic, []byte("online"))

	return &agent{
		pub: pub,

		groupedSensors: groupByInterval(sens),

		stateBase:         stateBase,
		discoveryBase:     s.DiscoveryPrefix,
		availabilityTopic: availabilityTopic,

		deviceId:     s.DeviceId,
		deviceName:   s.DeviceName,
		manufacturer: s.Manufacturer,
		model:        s.Model,

		once: s.Once,
	}, nil
}

func (a *agent) Run(ctx context.Context) error {
	if err := a.pub.Connect(10 * time.Second); err != nil {
		return err
	}
	defer func() {
		if a.availabilityTopic != "" {
			if err := a.pub.Publish(a.availabilityTopic, 1, true, []byte("offline")); err != nil {
				slog.Warn("mqtt publish offline failed", "topic", a.availabilityTopic, "err", err)
			}
		}
		a.pub.Close()
	}()

	if err := a.publishDiscovery(); err != nil {
		return err
	}

	if a.once {
		for _, group := range a.groupedSensors {
			sensorsStateCache := make(map[string]string, len(group))
			a.collectAndPublishGroup(ctx, group, sensorsStateCache)
		}
		return nil
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
	if err := a.pub.Connect(10 * time.Second); err != nil {
		return err
	}
	defer a.pub.Close()

	for _, group := range a.groupedSensors {
		for _, s := range group {
			topic := fmt.Sprintf("%s/sensor/%s/%s/config", a.discoveryBase, a.deviceId, s.Key())
			if err := a.pub.Publish(topic, 1, true, []byte{}); err != nil {
				return fmt.Errorf("purge: clear discovery failed (topic=%s): %w", topic, err)
			}
		}
	}

	for _, group := range a.groupedSensors {
		for _, s := range group {
			topic := fmt.Sprintf("%s/%s/state", a.stateBase, s.Key())
			if err := a.pub.Publish(topic, 1, true, []byte{}); err != nil {
				return fmt.Errorf("purge: clear state failed (topic=%s): %w", topic, err)
			}
		}
	}

	if a.availabilityTopic != "" {
		if err := a.pub.Publish(a.availabilityTopic, 1, true, []byte{}); err != nil {
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
