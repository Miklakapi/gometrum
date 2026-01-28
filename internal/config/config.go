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

	applyDefaults(&cfg)
	return cfg, nil
}

func ValidateConfig(cfg Config) error {
	if strings.TrimSpace(cfg.MQTT.Host) == "" {
		return errors.New("config: mqtt.host is required")
	}
	if cfg.MQTT.Port <= 0 || cfg.MQTT.Port > 65535 {
		return errors.New("config: mqtt.port must be a valid TCP port (1-65535)")
	}
	if strings.TrimSpace(cfg.MQTT.Username) != "" && strings.TrimSpace(cfg.MQTT.Password) == "" {
		return errors.New("config: mqtt.password is required when mqtt.username is set")
	}
	if strings.TrimSpace(cfg.MQTT.Password) != "" && strings.TrimSpace(cfg.MQTT.Username) == "" {
		return errors.New("config: mqtt.username is required when mqtt.password is set")
	}
	if strings.TrimSpace(cfg.MQTT.ClientID) == "" {
		return errors.New("config: mqtt.client_id is required (must be unique per device)")
	}
	if strings.TrimSpace(cfg.MQTT.DiscoveryPrefix) == "" {
		return errors.New("config: mqtt.discovery_prefix cannot be empty")
	}
	if strings.TrimSpace(cfg.MQTT.StatePrefix) == "" {
		return errors.New("config: mqtt.state_prefix cannot be empty")
	}
	if cfg.MQTT.DefaultInterval <= 0 {
		return errors.New("config: mqtt.default_interval must be > 0 (e.g. \"30s\")")
	}

	if strings.TrimSpace(cfg.Agent.DeviceID) == "" {
		return errors.New("config: agent.device_id is required (must be unique per device)")
	}
	if strings.TrimSpace(cfg.Agent.DeviceName) == "" {
		return errors.New("config: agent.device_name cannot be empty")
	}
	if strings.TrimSpace(cfg.Agent.Manufacturer) == "" {
		return errors.New("config: agent.manufacturer cannot be empty")
	}
	if strings.TrimSpace(cfg.Agent.Model) == "" {
		return errors.New("config: agent.model cannot be empty")
	}

	if len(cfg.Sensors) == 0 {
		return errors.New("config: sensors section must not be empty (define at least one sensor)")
	}

	for sensorKey, sensorCfg := range cfg.Sensors {
		if strings.TrimSpace(sensorKey) == "" {
			return errors.New("config: sensors contains an empty key")
		}
		if sensorCfg.Interval <= 0 {
			return errors.New("config: sensors." + sensorKey + ".interval resolved to 0 (check mqtt.default_interval)")
		}

		if len(sensorCfg.IncludeMounts) > 0 || len(sensorCfg.ExcludeMounts) > 0 {
			includeSet := make(map[string]struct{}, len(sensorCfg.IncludeMounts))
			for _, m := range sensorCfg.IncludeMounts {
				m = strings.TrimSpace(m)
				if m == "" {
					return errors.New("config: sensors." + sensorKey + ".include_mounts contains an empty mount")
				}
				includeSet[m] = struct{}{}
			}

			excludeSet := make(map[string]struct{}, len(sensorCfg.ExcludeMounts))
			for _, m := range sensorCfg.ExcludeMounts {
				m = strings.TrimSpace(m)
				if m == "" {
					return errors.New("config: sensors." + sensorKey + ".exclude_mounts contains an empty mount")
				}
				excludeSet[m] = struct{}{}
			}

			for m := range includeSet {
				if _, ok := excludeSet[m]; ok {
					return errors.New("config: sensors." + sensorKey + " mount is present in both include_mounts and exclude_mounts: " + m)
				}
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
