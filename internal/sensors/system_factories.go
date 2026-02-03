package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type uptimeSensor struct {
	base
}

func newUptimeSensor(key string, cfg config.SensorConfig) Sensor {
	return &uptimeSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *uptimeSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type osVersionSensor struct {
	base
}

func newOSVersionSensor(key string, cfg config.SensorConfig) Sensor {
	return &osVersionSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *osVersionSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
