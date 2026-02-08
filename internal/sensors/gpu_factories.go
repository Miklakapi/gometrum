package sensors

import (
	"context"
	"fmt"

	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
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
	st, err := nvmlStats()
	if err != nil {
		return "unavailable", nil
	}
	return fmt.Sprintf("%d", st.UtilGPU), nil
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
	st, err := nvmlStats()
	if err != nil {
		return "unavailable", nil
	}
	if st.MemTotal == 0 {
		return "0", nil
	}
	percent := (float64(st.MemUsed) / float64(st.MemTotal)) * 100.0
	return fmt.Sprintf("%.1f", percent), nil
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
	st, err := nvmlStats()
	if err != nil {
		return "unavailable", nil
	}
	return fmt.Sprintf("%d", st.TempC), nil
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
	st, err := nvmlStats()
	if err != nil {
		return "unavailable", nil
	}
	watts := float64(st.PowerMilliW) / 1000.0
	return fmt.Sprintf("%.1f", watts), nil
}

type gpuStats struct {
	UtilGPU     uint32
	MemUsed     uint64
	MemTotal    uint64
	TempC       uint32
	PowerMilliW uint32
}

func nvmlStats() (gpuStats, error) {
	var st gpuStats

	if ret := nvml.Init(); ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml init failed: %s", nvml.ErrorString(ret))
	}
	defer nvml.Shutdown()

	dev, ret := nvml.DeviceGetHandleByIndex(0)
	if ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml get device failed: %s", nvml.ErrorString(ret))
	}

	util, ret := dev.GetUtilizationRates()
	if ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml utilization failed: %s", nvml.ErrorString(ret))
	}

	mem, ret := dev.GetMemoryInfo()
	if ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml memory failed: %s", nvml.ErrorString(ret))
	}

	temp, ret := dev.GetTemperature(nvml.TEMPERATURE_GPU)
	if ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml temperature failed: %s", nvml.ErrorString(ret))
	}

	power, ret := dev.GetPowerUsage()
	if ret != nvml.SUCCESS {
		return st, fmt.Errorf("nvml power failed: %s", nvml.ErrorString(ret))
	}

	st.UtilGPU = util.Gpu
	st.MemUsed = mem.Used
	st.MemTotal = mem.Total
	st.TempC = temp
	st.PowerMilliW = power
	return st, nil
}
