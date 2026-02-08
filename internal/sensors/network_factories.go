package sensors

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

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
	_, level, err := readWiFiLevelFromProc()
	if err != nil {
		return "unavailable", err
	}

	if level > 0 {
		level = -level
	}
	return strconv.Itoa(level), nil
}

type wiFiSSIDSensor struct {
	base
}

func newWiFiSSIDSensor(key string, cfg config.SensorConfig) Sensor {
	return &wiFiSSIDSensor{
		base: base{key: key, name: cfg.Name, interval: cfg.Interval, ha: cfg.HA},
	}
}

func (s *wiFiSSIDSensor) Collect(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "iwgetid", "-r")
	out, err := cmd.Output()
	if err != nil {
		return "unavailable", err
	}

	ssid := strings.TrimSpace(string(out))
	if ssid == "" {
		return "unavailable", nil
	}
	return ssid, nil
}

func readWiFiLevelFromProc() (iface string, level int, err error) {
	data, err := os.ReadFile("/proc/net/wireless")
	if err != nil {
		return "", 0, fmt.Errorf("wifi_signal: read /proc/net/wireless failed: %w", err)
	}

	sc := bufio.NewScanner(bytes.NewReader(data))
	lineNo := 0
	for sc.Scan() {
		lineNo++
		if lineNo <= 2 {
			continue
		}

		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		iface = strings.TrimSuffix(parts[0], ":")
		levelStr := strings.TrimSuffix(parts[3], ".")

		f, ferr := strconv.ParseFloat(levelStr, 64)
		if ferr != nil {
			continue
		}

		level = int(f)
		return iface, level, nil
	}

	if err := sc.Err(); err != nil {
		return "", 0, fmt.Errorf("wifi_signal: scan /proc/net/wireless failed: %w", err)
	}

	return "", 0, fmt.Errorf("wifi_signal: no wifi interface in /proc/net/wireless")
}
