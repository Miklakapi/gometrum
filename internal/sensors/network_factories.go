package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type publicIPSensor struct {
	base
}

func newPublicIPSensor(key string, cfg config.SensorConfig) Sensor {
	return &publicIPSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *publicIPSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type wiFiSignalSensor struct {
	base
}

func newWiFiSignalSensor(key string, cfg config.SensorConfig) Sensor {
	return &wiFiSignalSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *wiFiSignalSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
