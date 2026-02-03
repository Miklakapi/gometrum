package sensors

import (
	"context"
	"fmt"
	"strings"

	"github.com/Miklakapi/gometrum/internal/config"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	gsensors "github.com/shirou/gopsutil/v4/sensors"
)

type cpuUsageSensor struct {
	base
}

func newCPUUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuUsageSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *cpuUsageSensor) Collect(ctx context.Context) (string, error) {
	vals, err := cpu.Percent(0, false)
	if err != nil {
		return "", err
	}
	if len(vals) == 0 {
		return "0.0", nil
	}
	return fmt.Sprintf("%.1f", vals[0]), nil
}

type cpuLoadSensor struct {
	base
}

func newCPULoadSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuLoadSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *cpuLoadSensor) Collect(ctx context.Context) (string, error) {
	avg, err := load.Avg()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", avg.Load1), nil
}

type cpuTempSensor struct {
	base
}

func newCPUTempSensor(key string, cfg config.SensorConfig) Sensor {
	return &cpuTempSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *cpuTempSensor) Collect(ctx context.Context) (string, error) {
	temps, err := gsensors.SensorsTemperatures()
	if err != nil {
		return "", err
	}
	if len(temps) == 0 {
		return "", fmt.Errorf("cpu_temp: no temperature sensors")
	}

	best := temps[0]
	for _, t := range temps {
		k := strings.ToLower(t.SensorKey)
		if strings.Contains(k, "cpu") ||
			strings.Contains(k, "core") ||
			strings.Contains(k, "package") ||
			strings.Contains(k, "soc") {
			best = t
			break
		}
	}

	return fmt.Sprintf("%.1f", best.Temperature), nil
}
