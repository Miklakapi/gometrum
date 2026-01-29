package sensors

var _ = map[string]SensorDefinition{
	// CPU
	"cpu_usage": {
		DefaultName: "CPU usage",
		DefaultIcon: "mdi:cpu-64-bit",
		DefaultUnit: "%",
	},
	"cpu_load": {
		DefaultName: "CPU load",
		DefaultIcon: "mdi:chart-line",
		DefaultUnit: "",
	},
	"cpu_temp": {
		DefaultName: "CPU temperature",
		DefaultIcon: "mdi:thermometer",
		DefaultUnit: "°C",
	},

	// System
	"uptime": {
		DefaultName: "System uptime",
		DefaultIcon: "mdi:timer-outline",
		DefaultUnit: "s",
	},
	"os_version": {
		DefaultName: "OS version",
		DefaultIcon: "mdi:linux",
		DefaultUnit: "",
	},

	// Memory
	"memory_usage": {
		DefaultName: "Memory usage",
		DefaultIcon: "mdi:memory",
		DefaultUnit: "%",
	},
	"swap_usage": {
		DefaultName: "Swap usage",
		DefaultIcon: "mdi:swap-horizontal",
		DefaultUnit: "%",
	},

	// Disk
	"disk_usage": {
		DefaultName: "Disk usage",
		DefaultIcon: "mdi:harddisk",
		DefaultUnit: "%",
	},

	// Network
	"public_ip": {
		DefaultName: "Public IP",
		DefaultIcon: "mdi:ip",
		DefaultUnit: "",
	},
	"wifi_signal": {
		DefaultName: "Wi-Fi signal",
		DefaultIcon: "mdi:wifi",
		DefaultUnit: "%",
	},

	// GPU
	"gpu_usage": {
		DefaultName: "GPU usage",
		DefaultIcon: "mdi:gpu",
		DefaultUnit: "%",
	},
	"gpu_memory_usage": {
		DefaultName: "GPU memory usage",
		DefaultIcon: "mdi:memory",
		DefaultUnit: "%",
	},
	"gpu_temp": {
		DefaultName: "GPU temperature",
		DefaultIcon: "mdi:thermometer",
		DefaultUnit: "°C",
	},
	"gpu_power": {
		DefaultName: "GPU power",
		DefaultIcon: "mdi:flash",
		DefaultUnit: "W",
	},
}
