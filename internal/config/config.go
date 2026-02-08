package config

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	normalizeConfig(&cfg)
	applyDefaults(&cfg)
	return cfg, nil
}

func ValidateConfig(cfg Config) error {
	if cfg.MQTT.Host == "" {
		return errors.New("config: mqtt.host is required")
	}
	if cfg.MQTT.Port <= 0 || cfg.MQTT.Port > 65535 {
		return errors.New("config: mqtt.port must be a valid TCP port (1-65535)")
	}
	if cfg.MQTT.Username != "" && cfg.MQTT.Password == "" {
		return errors.New("config: mqtt.password is required when mqtt.username is set")
	}
	if cfg.MQTT.Password != "" && cfg.MQTT.Username == "" {
		return errors.New("config: mqtt.username is required when mqtt.password is set")
	}
	if cfg.MQTT.ClientID == "" {
		return errors.New("config: mqtt.client_id is required (must be unique per device)")
	}
	if cfg.MQTT.DiscoveryPrefix == "" {
		return errors.New("config: mqtt.discovery_prefix cannot be empty")
	}
	if cfg.MQTT.StatePrefix == "" {
		return errors.New("config: mqtt.state_prefix cannot be empty")
	}
	if cfg.MQTT.DefaultInterval <= 0 {
		return errors.New("config: mqtt.default_interval must be > 0 (e.g. \"30s\")")
	}

	if cfg.Agent.DeviceID == "" {
		return errors.New("config: agent.device_id is required (must be unique per device)")
	}
	if cfg.Agent.DeviceName == "" {
		return errors.New("config: agent.device_name cannot be empty")
	}
	if cfg.Agent.Manufacturer == "" {
		return errors.New("config: agent.manufacturer cannot be empty")
	}
	if cfg.Agent.Model == "" {
		return errors.New("config: agent.model cannot be empty")
	}

	if len(cfg.Sensors) == 0 {
		return errors.New("config: sensors section must not be empty (define at least one sensor)")
	}

	for sensorKey, sensorCfg := range cfg.Sensors {
		if sensorKey == "" {
			return errors.New("config: sensors contains an empty key")
		}
		if sensorCfg.Interval <= 0 {
			return errors.New("config: sensors." + sensorKey + ".interval resolved to 0 (check mqtt.default_interval)")
		}

		if len(sensorCfg.IncludeMounts) > 0 {
			seen := make(map[string]struct{}, len(sensorCfg.IncludeMounts))

			for _, m := range sensorCfg.IncludeMounts {
				if m == "" {
					return errors.New("config: sensors." + sensorKey + ".include_mounts contains an empty mount")
				}
				if _, ok := seen[m]; ok {
					return errors.New("config: sensors." + sensorKey + ".include_mounts contains duplicate mount: " + m)
				}
				seen[m] = struct{}{}
			}
		}
	}

	return nil
}

func loadBytes(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func normalizeConfig(cfg *Config) {
	cfg.MQTT.Host = strings.TrimSpace(cfg.MQTT.Host)
	cfg.MQTT.Username = strings.TrimSpace(cfg.MQTT.Username)
	cfg.MQTT.Password = strings.TrimSpace(cfg.MQTT.Password)
	cfg.MQTT.ClientID = strings.TrimSpace(cfg.MQTT.ClientID)
	cfg.MQTT.DiscoveryPrefix = strings.TrimSpace(cfg.MQTT.DiscoveryPrefix)
	cfg.MQTT.StatePrefix = strings.TrimSpace(cfg.MQTT.StatePrefix)

	cfg.Agent.DeviceID = strings.TrimSpace(cfg.Agent.DeviceID)
	cfg.Agent.DeviceName = strings.TrimSpace(cfg.Agent.DeviceName)
	cfg.Agent.Manufacturer = strings.TrimSpace(cfg.Agent.Manufacturer)
	cfg.Agent.Model = strings.TrimSpace(cfg.Agent.Model)

	normalized := make(map[string]SensorConfig, len(cfg.Sensors))

	for key, sensor := range cfg.Sensors {
		sensorKey := strings.TrimSpace(key)
		if sensorKey == "" {
			continue
		}

		sensor.Name = strings.TrimSpace(sensor.Name)

		for i, m := range sensor.IncludeMounts {
			sensor.IncludeMounts[i] = strings.TrimSpace(m)
		}

		if sensor.HA != nil {
			sensor.HA.Icon = strings.TrimSpace(sensor.HA.Icon)
			sensor.HA.Unit = strings.TrimSpace(sensor.HA.Unit)
		}

		normalized[sensorKey] = sensor
	}

	cfg.Sensors = normalized
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
	if cfg.MQTT.DefaultInterval <= 0 {
		cfg.MQTT.DefaultInterval = time.Duration(1 * time.Minute)
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

	for key, sensorCfg := range cfg.Sensors {
		if sensorCfg.Interval <= 0 {
			sensorCfg.Interval = cfg.MQTT.DefaultInterval
			cfg.Sensors[key] = sensorCfg
		}
	}
}
