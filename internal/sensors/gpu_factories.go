package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type gpuUsageSensor struct {
	base
}

func newGPUUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &gpuUsageSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *gpuUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type gpuMemoryUsageSensor struct {
	base
}

func newGPUMemoryUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &gpuMemoryUsageSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *gpuMemoryUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type gpuTempSensor struct {
	base
}

func newGPUTempSensor(key string, cfg config.SensorConfig) Sensor {
	return &gpuTempSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *gpuTempSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}

type gpuPowerSensor struct {
	base
}

func newGPUPowerSensor(key string, cfg config.SensorConfig) Sensor {
	return &gpuPowerSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *gpuPowerSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
