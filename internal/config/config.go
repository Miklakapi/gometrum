package config

import (
	_ "embed"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

//go:embed example.yaml
var ExampleYAML string

func SaveExample(path string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(ExampleYAML), 0644)
}

func LoadString(path string) (string, error) {
	data, err := loadBytes(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func LoadAndValidate(path string) (Config, error) {
	cfg, err := LoadConfig(path)
	if err != nil {
		return cfg, err
	}

	if err := ValidateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	data, err := loadBytes(path)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	applyDefaults(&cfg)
	return cfg, nil
}

func ValidateConfig(cfg Config) error {
	panic("TODO")
}

func loadBytes(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func applyDefaults(cfg *Config) {
	if cfg.MQTT.Port == 0 {
		cfg.MQTT.Port = 1883
	}
	if cfg.MQTT.DiscoveryPrefix == "" {
		cfg.MQTT.DiscoveryPrefix = "homeassistant"
	}
	if cfg.MQTT.StatePrefix == "" {
		cfg.MQTT.StatePrefix = "gometrum"
	}

	if cfg.Agent.DeviceName == "" {
		cfg.Agent.DeviceName = "gometrum"
	}
	if cfg.Agent.Manufacturer == "" {
		cfg.Agent.Manufacturer = "gometrum"
	}
	if cfg.Agent.Model == "" {
		cfg.Agent.Model = "linux-host"
	}
}
