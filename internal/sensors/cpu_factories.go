package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type cpuUsageSensor struct {
	base
}

func newCPUUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuUsageSensor{
		base: base{
			key:      key,
			name:     cfg.Name,
			interval: cfg.Interval,
			ha:       cfg.HA,
		},
	}
}

func (s *cpuUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type cpuLoadSensor struct {
	base
}

func newCPULoadSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuLoadSensor{
		base: base{
			key:      key,
			name:     cfg.Name,
			interval: cfg.Interval,
			ha:       cfg.HA,
		},
	}
}

func (s *cpuLoadSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type cpuTempSensor struct {
	base
}

func newCPUTempSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuTempSensor{
		base: base{
			key:      key,
			name:     cfg.Name,
			interval: cfg.Interval,
			ha:       cfg.HA,
		},
	}
}

func (s *cpuTempSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
