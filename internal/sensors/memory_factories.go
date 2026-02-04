package sensors

import (
	"context"
	"fmt"

	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/shirou/gopsutil/v4/mem"
)

type memoryUsageSensor struct {
	base
}

func newMemoryUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &memoryUsageSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *memoryUsageSensor) Collect(ctx context.Context) (string, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return "unavailable", fmt.Errorf("memory_usage: %w", err)
	}

	return fmt.Sprintf("%.1f", vm.UsedPercent), nil
}

type swapUsageSensor struct {
	base
}

func newSwapUsageSensor(key string, cfg config.SensorConfig) Sensor {
	return &swapUsageSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *swapUsageSensor) Collect(ctx context.Context) (string, error) {
	sm, err := mem.SwapMemory()
	if err != nil {
		return "unavailable", fmt.Errorf("swap_usage: %w", err)
	}

	if sm.Total == 0 {
		return "0.0", nil
	}

	return fmt.Sprintf("%.1f", sm.UsedPercent), nil
}
