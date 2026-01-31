package sensors

import (
	"context"
	"fmt"
	"time"

	"github.com/Miklakapi/gometrum/internal/config"
)

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

type Sensor interface {
	Key() string
	Name() string
	Interval() time.Duration
	HA() *config.HASensorConfig
	Collect(ctx context.Context) (string, error)
}

func Build(cfg config.Config) ([]Sensor, error) {
	out := make([]Sensor, 0, len(cfg.Sensors))

	for key, scfg := range cfg.Sensors {
		f, ok := factories[key]
		if !ok {
			return nil, fmt.Errorf("sensors: no factory for %s", key)
		}
		s, err := f(key, scfg)
		if err != nil {
			return nil, fmt.Errorf("sensors: %s: %w", key, err)
		}
		out = append(out, s)
	}

	return out, nil
}

var factories = map[string]func(key string, cfg config.SensorConfig) (Sensor, error){
	"uptime": func(key string, cfg config.SensorConfig) (Sensor, error) {
		return newUptimeSensor(key, cfg), nil
	},
}
