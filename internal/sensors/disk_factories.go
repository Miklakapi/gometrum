package sensors

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Miklakapi/gometrum/internal/config"
	"github.com/shirou/gopsutil/v4/disk"
)

type diskUsageSensor struct {
	base
	mount string
}

func newDiskUsageSensors(key string, cfg config.SensorConfig) []Sensor {
	include := cfg.IncludeMounts

	if len(include) == 0 {
		return nil
	}

	sort.Strings(include)

	out := make([]Sensor, 0, len(include))
	for _, m := range include {
		sKey := key + "_" + sanitizeMount(m)

		sName := cfg.Name
		if m != "/" {
			sName = fmt.Sprintf("%s %s", cfg.Name, m)
		}

		out = append(out, &diskUsageSensor{
			base:  base{key: sKey, name: sName, interval: cfg.Interval, ha: cfg.HA},
			mount: m,
		})
	}

	return out
}

func (s *diskUsageSensor) Collect(ctx context.Context) (string, error) {
	u, err := disk.UsageWithContext(ctx, s.mount)
	if err != nil {
		return "unavailable", fmt.Errorf("disk_usage(%s): %w", s.mount, err)
	}
	return fmt.Sprintf("%.1f", u.UsedPercent), nil
}

func sanitizeMount(m string) string {
	if m == "/" {
		return "root"
	}
	s := strings.Trim(m, "/")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "-", "_")
	if s == "" {
		return "mount"
	}
	return s
}
