package sensors

import (
	"context"
	"fmt"
	"net"

	"github.com/Miklakapi/gometrum/internal/config"
)

type hostIPSensor struct {
	base
}

func newHostIPSensor(key string, cfg config.SensorConfig) Sensor {
	return &hostIPSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *hostIPSensor) Collect(ctx context.Context) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unavailable", fmt.Errorf("host_ip: %w", err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		ip := ipNet.IP
		if ip == nil || ip.IsLoopback() {
			continue
		}

		ip4 := ip.To4()
		if ip4 == nil {
			continue
		}

		return ip4.String(), nil
	}

	return "unavailable", fmt.Errorf("host_ip: no non-loopback IPv4 found")
}

type wiFiSignalSensor struct {
	base
}

func newWiFiSignalSensor(key string, cfg config.SensorConfig) Sensor {
	return &wiFiSignalSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *wiFiSignalSensor) Collect(ctx context.Context) (string, error) {
	return "test", nil
}
