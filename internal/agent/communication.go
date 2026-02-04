package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Miklakapi/gometrum/internal/sensors"
)

type haSensorDiscovery struct {
	Name              string    `json:"name"`
	UniqueID          string    `json:"unique_id"`
	StateTopic        string    `json:"state_topic"`
	AvailabilityTopic string    `json:"availability_topic,omitempty"`
	Icon              string    `json:"icon,omitempty"`
	Unit              string    `json:"unit_of_measurement,omitempty"`
	Device            *haDevice `json:"device,omitempty"`
}

type haDevice struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
}

func (a *agent) publishDiscovery() error {
	dev := &haDevice{
		Identifiers:  []string{a.deviceId},
		Name:         a.deviceName,
		Manufacturer: a.manufacturer,
		Model:        a.model,
	}

	for _, group := range a.groupedSensors {
		for _, s := range group {
			key := s.Key()

			stateTopic := fmt.Sprintf("%s/%s/state", a.stateBase, key)
			configTopic := fmt.Sprintf("%s/sensor/%s/%s/config", a.discoveryBase, a.deviceId, key)

			payload := haSensorDiscovery{
				Name:              s.Name(),
				UniqueID:          fmt.Sprintf("%s_%s", a.deviceId, key),
				StateTopic:        stateTopic,
				AvailabilityTopic: a.availabilityTopic,
				Device:            dev,
			}

			if ha := s.HA(); ha != nil {
				if ha.Icon != "" {
					payload.Icon = ha.Icon
				}
				if ha.Unit != "" {
					payload.Unit = ha.Unit
				}
			}

			b, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("discovery marshal failed (sensor=%s): %w", key, err)
			}

			if err := a.client.Publish(configTopic, 1, true, b); err != nil {
				return fmt.Errorf("discovery publish failed (sensor=%s, topic=%s): %w", key, configTopic, err)
			}
		}
	}

	return nil
}

func (a *agent) collectAndPublishGroup(ctx context.Context, group []sensors.Sensor, sensorsStateCache map[string]string) {
	for _, s := range group {
		topic := fmt.Sprintf("%s/%s/state", a.stateBase, s.Key())

		val, err := s.Collect(ctx)
		if err != nil {
			slog.Error("collect failed", "sensor", s.Key(), "err", err)
		}

		if prev, ok := sensorsStateCache[s.Key()]; ok && prev == val {
			continue
		}
		sensorsStateCache[s.Key()] = val

		if err := a.client.Publish(topic, 1, true, []byte(val)); err != nil {
			slog.Error("publish failed", "sensor", s.Key(), "topic", topic, "err", err)
		} else {
			slog.Info("published", "sensor", s.Key(), "topic", topic, "value", val)
		}
	}
}
