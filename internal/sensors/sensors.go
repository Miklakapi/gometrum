package sensors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Miklakapi/gometrum/internal/config"
)

type Sensor interface {
	Key() string
	Name() string
	Interval() time.Duration
	HA() *config.HASensorConfig
	Collect(ctx context.Context) (string, error)
}

type base struct {
	key      string
	name     string
	interval time.Duration
	ha       *config.HASensorConfig
}

func (b base) Key() string                { return b.key }
func (b base) Name() string               { return b.name }
func (b base) Interval() time.Duration    { return b.interval }
func (b base) HA() *config.HASensorConfig { return b.ha }

type SensorDefinition struct {
	DefaultName string
	DefaultIcon string
	DefaultUnit string
	Factory     func(key string, cfg config.SensorConfig) ([]Sensor, error)
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

		if sensorCfg.Name == "" {
			sensorCfg.Name = def.DefaultName
		}

		if sensorCfg.Interval <= 0 {
			sensorCfg.Interval = cfg.MQTT.DefaultInterval
		}

		if def.DefaultIcon != "" || def.DefaultUnit != "" {
			if sensorCfg.HA == nil {
				sensorCfg.HA = &config.HASensorConfig{}
			}
			if sensorCfg.HA.Icon == "" && def.DefaultIcon != "" {
				sensorCfg.HA.Icon = def.DefaultIcon
			}
			if sensorCfg.HA.Unit == "" && def.DefaultUnit != "" {
				sensorCfg.HA.Unit = def.DefaultUnit
			}
		}

		cfg.Sensors[sensorKey] = sensorCfg
	}
	return nil
}

func Validate(cfg config.Config) error {
	for sensorKey, sensorCfg := range cfg.Sensors {
		def, ok := registry[sensorKey]
		if !ok {
			return errors.New("sensors: unknown sensor type: " + sensorKey)
		}
		if def.Factory == nil {
			return errors.New("sensors." + sensorKey + ": no Factory implementation")
		}

		if sensorCfg.Name == "" {
			return errors.New("sensors." + sensorKey + ": name is empty and no DefaultName in registry")
		}
		if sensorCfg.Interval <= 0 {
			return errors.New("sensors." + sensorKey + ": interval resolved to 0 (check mqtt.default_interval)")
		}
	}

	return nil
}

func Build(cfg config.Config) ([]Sensor, error) {
	out := make([]Sensor, 0, len(cfg.Sensors))

	for key, scfg := range cfg.Sensors {
		def, ok := registry[key]
		if !ok {
			return nil, fmt.Errorf("sensors: unknown sensor type: %s", key)
		}
		if def.Factory == nil {
			return nil, fmt.Errorf("sensors: %s has no Factory (not implemented)", key)
		}

		list, err := def.Factory(key, scfg)
		if err != nil {
			return nil, fmt.Errorf("sensors: %s: %w", key, err)
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("sensors: %s factory returned 0 sensors", key)
		}

		out = append(out, list...)
	}

	return out, nil
}
