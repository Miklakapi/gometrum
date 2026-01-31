package sensors

import (
	"context"

	"github.com/Miklakapi/gometrum/internal/config"
)

type diskUsageSensor struct {
	base
}

func newDiskUsageSensors(key string, cfg config.SensorConfig) []Sensor {
	return []Sensor{
		&uptimeSensor{
			base: base{
				key:      key,
				name:     cfg.Name,
				interval: cfg.Interval,
				ha:       cfg.HA,
			},
		},
	}
}

func (s *diskUsageSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
