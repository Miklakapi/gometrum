package sensors

import (
	"github.com/Miklakapi/gometrum/internal/config"
)

var registry = map[string]SensorDefinition{
	// CPU
	"cpu_usage": {
		DefaultName: "CPU usage",
		DefaultIcon: "mdi:cpu-64-bit",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newCPUUsageSensor(key, cfg)}, nil
		},
	},
	"cpu_load_1m": {
		DefaultName: "CPU load (1m)",
		DefaultIcon: "mdi:chart-line",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newCPULoadSensor(key, cfg, loadWindow1m)}, nil
		},
	},
	"cpu_load_5m": {
		DefaultName: "CPU load (5m)",
		DefaultIcon: "mdi:chart-line",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newCPULoadSensor(key, cfg, loadWindow5m)}, nil
		},
	},
	"cpu_load_15m": {
		DefaultName: "CPU load (15m)",
		DefaultIcon: "mdi:chart-line",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newCPULoadSensor(key, cfg, loadWindow15m)}, nil
		},
	},
	"cpu_temp": {
		DefaultName: "CPU temperature",
		DefaultIcon: "mdi:thermometer",
		DefaultUnit: "°C",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newCPUTempSensor(key, cfg)}, nil
		},
	},

	// System
	"uptime": {
		DefaultName: "System uptime",
		DefaultIcon: "mdi:timer-outline",
		DefaultUnit: "s",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newUptimeSensor(key, cfg)}, nil
		},
	},
	"os_version": {
		DefaultName: "OS version",
		DefaultIcon: "mdi:linux",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newOSVersionSensor(key, cfg)}, nil
		},
	},
	"hostname": {
		DefaultName: "Hostname",
		DefaultIcon: "mdi:server",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newHostnameSensor(key, cfg)}, nil
		},
	},

	// Memory
	"memory_usage": {
		DefaultName: "Memory usage",
		DefaultIcon: "mdi:memory",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newMemoryUsageSensor(key, cfg)}, nil
		},
	},
	"swap_usage": {
		DefaultName: "Swap usage",
		DefaultIcon: "mdi:swap-horizontal",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newSwapUsageSensor(key, cfg)}, nil
		},
	},

	// Disk
	"disk_usage": {
		DefaultName: "Disk usage",
		DefaultIcon: "mdi:harddisk",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return newDiskUsageSensors(key, cfg), nil
		},
	},

	// Network
	"host_ip": {
		DefaultName: "Public IP",
		DefaultIcon: "mdi:ip",
		DefaultUnit: "",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newHostIPSensor(key, cfg)}, nil
		},
	},
	"wifi_signal": {
		DefaultName: "Wi-Fi signal",
		DefaultIcon: "mdi:wifi",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newWiFiSignalSensor(key, cfg)}, nil
		},
	},

	// GPU
	"gpu_usage": {
		DefaultName: "GPU usage",
		DefaultIcon: "mdi:gpu",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newGPUUsageSensor(key, cfg)}, nil
		},
	},
	"gpu_memory_usage": {
		DefaultName: "GPU memory usage",
		DefaultIcon: "mdi:memory",
		DefaultUnit: "%",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newGPUMemoryUsageSensor(key, cfg)}, nil
		},
	},
	"gpu_temp": {
		DefaultName: "GPU temperature",
		DefaultIcon: "mdi:thermometer",
		DefaultUnit: "°C",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newGPUTempSensor(key, cfg)}, nil
		},
	},
	"gpu_power": {
		DefaultName: "GPU power",
		DefaultIcon: "mdi:flash",
		DefaultUnit: "W",
		Factory: func(key string, cfg config.SensorConfig) ([]Sensor, error) {
			return []Sensor{newGPUPowerSensor(key, cfg)}, nil
		},
	},
}
