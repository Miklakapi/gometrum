package sensors

import (
	"errors"
	"strings"

	"github.com/Miklakapi/gometrum/internal/config"
)

type SensorDefinition struct {
	DefaultName string
	DefaultIcon string
	DefaultUnit string
}

type Sensor struct {
}

func Prepare(cfg *config.Config) error {
	if err := Normalize(cfg); err != nil {
		return err
	}
	if err := Validate(*cfg); err != nil {
		return err
	}
	return nil
}

func Normalize(cfg *config.Config) error {
	for sensorKey, sensorCfg := range cfg.Sensors {
		def, ok := registry[sensorKey]
		if !ok {
			return errors.New("sensors: unknown sensor type: " + sensorKey)
		}

		if strings.TrimSpace(sensorCfg.Name) == "" {
			sensorCfg.Name = def.DefaultName
		}

		if sensorCfg.Interval <= 0 {
			sensorCfg.Interval = cfg.MQTT.DefaultInterval
		}

		if def.DefaultIcon != "" || def.DefaultUnit != "" {
			if sensorCfg.HA == nil {
				sensorCfg.HA = &config.HASensorConfig{}
			}
			if strings.TrimSpace(sensorCfg.HA.Icon) == "" && def.DefaultIcon != "" {
				sensorCfg.HA.Icon = def.DefaultIcon
			}
			if strings.TrimSpace(sensorCfg.HA.Unit) == "" && def.DefaultUnit != "" {
				sensorCfg.HA.Unit = def.DefaultUnit
			}
		}

		cfg.Sensors[sensorKey] = sensorCfg
	}
	return nil
}

func Validate(cfg config.Config) error {
	for sensorKey, sensorCfg := range cfg.Sensors {
		_, ok := registry[sensorKey]
		if !ok {
			return errors.New("sensors: unknown sensor type: " + sensorKey)
		}

		if strings.TrimSpace(sensorCfg.Name) == "" {
			return errors.New("sensors." + sensorKey + ": name is empty and no DefaultName in registry")
		}
		if sensorCfg.Interval <= 0 {
			return errors.New("sensors." + sensorKey + ": interval resolved to 0 (check mqtt.default_interval)")
		}
	}

	return nil
}

func Build(cfg config.Config) ([]Sensor, error) {
	panic("TODO")
}
