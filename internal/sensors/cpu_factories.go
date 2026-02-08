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

type loadWindow uint8

const (
	loadWindow1m loadWindow = iota
	loadWindow5m
	loadWindow15m
)

type cpuLoadSensor struct {
	base
	window loadWindow
}

func newCPULoadSensor(key string, cfg config.SensorConfig, window loadWindow) Sensor {
	return &cpuLoadSensor{
		base:   base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
		window: window,
	}
}

func (s *cpuLoadSensor) Collect(ctx context.Context) (string, error) {
	avg, err := load.AvgWithContext(ctx)
	if err != nil {
		return "unavailable", err
	}

	switch s.window {
	case loadWindow1m:
		return fmt.Sprintf("%.2f", avg.Load1), nil
	case loadWindow5m:
		return fmt.Sprintf("%.2f", avg.Load5), nil
	case loadWindow15m:
		return fmt.Sprintf("%.2f", avg.Load15), nil
	default:
		return "unavailable", fmt.Errorf("cpu_load: unknown window")
	}
}

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
		return "unavailable", err
	}
	if len(vals) == 0 {
		return "0.0", nil
	}
	return fmt.Sprintf("%.1f", vals[0]), nil
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
		return "unavailable", err
	}
	if len(temps) == 0 {
		return "unavailable", fmt.Errorf("cpu_temp: no temperature sensors")
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
