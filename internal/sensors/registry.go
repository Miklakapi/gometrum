package sensors

var Sensor = map[string]SensorDefinition{
	"cpu_usage": {
		DefaultName: "CPU usage",
		DefaultIcon: "",
		DefaultUnit: "",
	},
	"memory_usage": {
		DefaultName: "RAM usage",
		DefaultIcon: "",
		DefaultUnit: "",
	},
	"disk_usage": {
		DefaultName: "Disk usage",
		DefaultIcon: "",
		DefaultUnit: "",
	},
	"wifi_signal": {
		DefaultName: "WiFi signal",
		DefaultIcon: "",
		DefaultUnit: "",
	},
}
