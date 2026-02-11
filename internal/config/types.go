package config

import "time"

type Config struct {
	Log     LogConfig               `yaml:"log"`
	MQTT    MQTTConfig              `yaml:"mqtt"`
	Agent   AgentConfig             `yaml:"agent"`
	Sensors map[string]SensorConfig `yaml:"sensors"`
}

type LogConfig struct {
	Level string    `yaml:"level"`
	Sinks []LogSink `yaml:"sinks"`
}

type LogSink struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`

	Level string `yaml:"level"`

	Addr string `yaml:"addr"`

	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Timeout time.Duration     `yaml:"timeout"`
	Headers map[string]string `yaml:"headers"`

	Codec string `yaml:"codec"`
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
	HA            *HASensorConfig `yaml:"ha,omitempty"`
}

type HASensorConfig struct {
	Icon        string `yaml:"icon,omitempty"`
	Unit        string `yaml:"unit,omitempty"`
	DeviceClass string `yaml:"device_class,omitempty"`
	StateClass  string `yaml:"state_class,omitempty"`
}
