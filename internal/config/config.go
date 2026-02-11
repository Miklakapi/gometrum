package config

import (
	_ "embed"
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
	var err error

	if err = validateLogLevel(cfg.Log); err != nil {
		return err
	}

	if err = validateMQTT(cfg.MQTT); err != nil {
		return err
	}

	if err = validateAgent(cfg.Agent); err != nil {
		return err
	}

	if err = validateSensors(cfg.Sensors); err != nil {
		return err
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
	cfg.Log.Level = strings.ToLower(strings.TrimSpace(cfg.Log.Level))
	for i := range cfg.Log.Sinks {
		cfg.Log.Sinks[i].Type = strings.ToLower(strings.TrimSpace(cfg.Log.Sinks[i].Type))
		cfg.Log.Sinks[i].Name = strings.TrimSpace(cfg.Log.Sinks[i].Name)
		cfg.Log.Sinks[i].Level = strings.ToLower(strings.TrimSpace(cfg.Log.Sinks[i].Level))

		cfg.Log.Sinks[i].Addr = strings.TrimSpace(cfg.Log.Sinks[i].Addr)

		cfg.Log.Sinks[i].URL = strings.TrimSpace(cfg.Log.Sinks[i].URL)
		cfg.Log.Sinks[i].Method = strings.ToUpper(strings.TrimSpace(cfg.Log.Sinks[i].Method))
		cfg.Log.Sinks[i].Codec = strings.ToLower(strings.TrimSpace(cfg.Log.Sinks[i].Codec))

		if cfg.Log.Sinks[i].Headers != nil {
			normalized := make(map[string]string, len(cfg.Log.Sinks[i].Headers))
			for k, v := range cfg.Log.Sinks[i].Headers {
				key := strings.TrimSpace(k)
				val := strings.TrimSpace(v)
				if key == "" {
					continue
				}
				normalized[key] = val
			}
			cfg.Log.Sinks[i].Headers = normalized
		}
	}

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

	normalizedSensors := make(map[string]SensorConfig, len(cfg.Sensors))
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
			sensor.HA.DeviceClass = strings.TrimSpace(sensor.HA.DeviceClass)
			sensor.HA.StateClass = strings.TrimSpace(sensor.HA.StateClass)
		}

		normalizedSensors[sensorKey] = sensor
	}

	cfg.Sensors = normalizedSensors
}

func applyDefaults(cfg *Config) {
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	for i := range cfg.Log.Sinks {
		if cfg.Log.Sinks[i].Level == "" {
			cfg.Log.Sinks[i].Level = cfg.Log.Level
		}

		if cfg.Log.Sinks[i].Type == "http" {
			if cfg.Log.Sinks[i].Method == "" {
				cfg.Log.Sinks[i].Method = "POST"
			}
			if cfg.Log.Sinks[i].Timeout <= 0 {
				cfg.Log.Sinks[i].Timeout = 2 * time.Second
			}
			if cfg.Log.Sinks[i].Codec == "" {
				cfg.Log.Sinks[i].Codec = "event_json"
			}
			if cfg.Log.Sinks[i].Headers == nil {
				cfg.Log.Sinks[i].Headers = map[string]string{}
			}
		}
	}

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
