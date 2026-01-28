package config

import "time"

type Config struct {
	MQTT    MQTTConfig              `yaml:"mqtt"`
	Agent   AgentConfig             `yaml:"agent"`
	Sensors map[string]SensorConfig `yaml:"sensors"`
}

type MQTTConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	ClientID        string        `yaml:"client_id"`
	DiscoveryPrefix string        `yaml:"discovery_prefix"`
	StatePrefix     string        `yaml:"state_prefix"`
	DefaultInterval time.Duration `yaml:"default_interval"`
}

type AgentConfig struct {
	DeviceID     string `yaml:"device_id"`
	DeviceName   string `yaml:"device_name"`
	Manufacturer string `yaml:"manufacturer"`
	Model        string `yaml:"model"`
}

type SensorConfig struct {
	Name          string          `yaml:"name"`
	Interval      time.Duration   `yaml:"interval"`
	IncludeMounts []string        `yaml:"include_mounts,omitempty"`
	ExcludeMounts []string        `yaml:"exclude_mounts,omitempty"`
	HA            *HASensorConfig `yaml:"ha,omitempty"`
}

type HASensorConfig struct {
	Icon string `yaml:"icon,omitempty"`
	Unit string `yaml:"unit,omitempty"`
}
