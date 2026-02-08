package sensors

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/shirou/gopsutil/v4/host"
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
	u, err := host.Uptime()
	if err != nil {
		return "unavailable", err
	}
	return fmt.Sprintf("%d", u), nil
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
	hi, err := host.Info()
	if err != nil {
		return "unavailable", err
	}

	if hi.Platform != "" && hi.PlatformVersion != "" {
		return fmt.Sprintf("%s %s", hi.Platform, hi.PlatformVersion), nil
	}
	if hi.OS != "" && hi.KernelVersion != "" {
		return fmt.Sprintf("%s %s", hi.OS, hi.KernelVersion), nil
	}

	if hi.KernelVersion != "" {
		return hi.KernelVersion, nil
	}

	return "unavailable", errors.New("empty host info")
}

type hostnameSensor struct {
	base
}

func newHostnameSensor(key string, cfg config.SensorConfig) Sensor {
	return &hostnameSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *hostnameSensor) Collect(ctx context.Context) (string, error) {
	h, err := os.Hostname()
	if err != nil {
		return "unavailable", err
	}
	if h == "" {
		return "unavailable", errors.New("empty hostname")
	}
	return h, nil
}
