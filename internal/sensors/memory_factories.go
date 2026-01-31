package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type memoryUsageSensor struct {
	base
}

func newMemoryUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &memoryUsageSensor{
		base: base{
			key:      key,
			name:     cfg.Name,
			interval: cfg.Interval,
			ha:       cfg.HA,
		},
	}
}

func (s *memoryUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type swapUsageSensor struct {
	base
}

func newSwapUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &swapUsageSensor{
		base: base{
			key:      key,
			name:     cfg.Name,
			interval: cfg.Interval,
			ha:       cfg.HA,
		},
	}
}

func (s *swapUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
