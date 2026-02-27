package config

import (
	"errors"
	"fmt"
	"net"
	neturl "net/url"
)

func validateLogLevel(lc LogConfig) error {
	if lc.Level == "" {
		return errors.New("config: log.level cannot be empty")
	}
	if !isValidLogLevel(lc.Level) {
		return fmt.Errorf("config: log.level must be one of: debug, info, warn, error (got: %s)", lc.Level)
	}

	for i, sink := range lc.Sinks {
		path := fmt.Sprintf("config: log.sinks[%d]", i)

		if sink.Type == "" {
			return errors.New(path + ".type is required")
		}
		if sink.Name == "" {
			return errors.New(path + ".name is required")
		}
		if sink.Level == "" {
			return errors.New(path + ".level cannot be empty")
		}
		if !isValidLogLevel(sink.Level) {
			return fmt.Errorf("%s.level must be one of: debug, info, warn, error (got: %s)", path, sink.Level)
		}

		if sink.QueueSize < 1 {
			return fmt.Errorf("%s.queue_size must be >= 1 (got: %d)", path, sink.QueueSize)
		}

		switch sink.Type {
		case "udp":
			if sink.Addr == "" {
				return errors.New(path + ".addr is required for type=udp")
			}
			if err := validateHostPort(sink.Addr); err != nil {
				return fmt.Errorf("%s.addr must be host:port (got: %s): %w", path, sink.Addr, err)
			}
			if sink.Codec != "" {
				return errors.New(path + ".codec is not supported for type=udp (udp is text/syslog-like only)")
			}
			if sink.Batch != nil {
				return errors.New(path + ".batch is not supported for type=udp")
			}
			if sink.URL != "" || sink.Method != "" || sink.Timeout != 0 || len(sink.Headers) > 0 {
				return errors.New(path + " contains http-only fields but type=udp")
			}

		case "http":
			if sink.URL == "" {
				return errors.New(path + ".url is required for type=http")
			}
			if err := validateHTTPURL(sink.URL); err != nil {
				return fmt.Errorf("%s.url is invalid: %w", path, err)
			}

			if sink.Method == "" {
				return errors.New(path + ".method cannot be empty for type=http")
			}
			if !isValidHTTPMethod(sink.Method) {
				return fmt.Errorf("%s.method must be one of: GET, POST, PUT, PATCH, DELETE (got: %s)", path, sink.Method)
			}

			if sink.Timeout <= 0 {
				return errors.New(path + ".timeout must be > 0 for type=http (e.g. \"2s\")")
			}

			if sink.Codec == "" {
				return errors.New(path + ".codec cannot be empty for type=http")
			}
			if !isValidHTTPCodec(sink.Codec) {
				return fmt.Errorf("%s.codec must be one of: event_json, text, loki, ndjson (got: %s)", path, sink.Codec)
			}

			if sink.Addr != "" {
				return errors.New(path + ".addr is not supported for type=http")
			}

			if sink.Batch != nil {
				if sink.Codec != "ndjson" && sink.Codec != "loki" {
					return fmt.Errorf("%s.batch is only supported for codec=ndjson or codec=loki (got: %s)", path, sink.Codec)
				}
				if sink.Batch.MaxItems < 1 {
					return fmt.Errorf("%s.batch.max_items must be >= 1 (got: %d)", path, sink.Batch.MaxItems)
				}
				if sink.Batch.MaxWait <= 0 {
					return fmt.Errorf("%s.batch.max_wait must be > 0 (got: %s)", path, sink.Batch.MaxWait)
				}
			}

		default:
			return fmt.Errorf("%s.type must be one of: udp, http (got: %s)", path, sink.Type)
		}
	}

	return nil
}

func validateMQTT(mc MQTTConfig) error {
	if mc.Host == "" {
		return errors.New("config: mqtt.host is required")
	}
	if mc.Port <= 0 || mc.Port > 65535 {
		return errors.New("config: mqtt.port must be a valid TCP port (1-65535)")
	}
	if mc.Username != "" && mc.Password == "" {
		return errors.New("config: mqtt.password is required when mqtt.username is set")
	}
	if mc.Password != "" && mc.Username == "" {
		return errors.New("config: mqtt.username is required when mqtt.password is set")
	}
	if mc.ClientID == "" {
		return errors.New("config: mqtt.client_id is required (must be unique per device)")
	}
	if mc.DiscoveryPrefix == "" {
		return errors.New("config: mqtt.discovery_prefix cannot be empty")
	}
	if mc.StatePrefix == "" {
		return errors.New("config: mqtt.state_prefix cannot be empty")
	}
	if mc.DefaultInterval <= 0 {
		return errors.New("config: mqtt.default_interval must be > 0 (e.g. \"30s\")")
	}

	return nil
}

func validateAgent(ac AgentConfig) error {
	if ac.DeviceID == "" {
		return errors.New("config: agent.device_id is required (must be unique per device)")
	}
	if ac.DeviceName == "" {
		return errors.New("config: agent.device_name cannot be empty")
	}
	if ac.Manufacturer == "" {
		return errors.New("config: agent.manufacturer cannot be empty")
	}
	if ac.Model == "" {
		return errors.New("config: agent.model cannot be empty")
	}

	return nil
}

func validateSensors(sc map[string]SensorConfig) error {
	if len(sc) == 0 {
		return errors.New("config: sensors section must not be empty (define at least one sensor)")
	}

	for sensorKey, sensorCfg := range sc {
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

func validateButtons(bc map[string]ButtonConfig) error {
	if len(bc) == 0 {
		return nil
	}

	for buttonKey, buttonCfg := range bc {
		if buttonKey == "" {
			return errors.New("config: buttons contains an empty key")
		}

		if buttonCfg.Name == "" {
			return errors.New("config: buttons." + buttonKey + ".name is required")
		}

		if len(buttonCfg.Command) == 0 {
			return errors.New("config: buttons." + buttonKey + ".command must contain at least one item (executable name)")
		}

		for i, arg := range buttonCfg.Command {
			if arg == "" {
				return fmt.Errorf("config: buttons.%s.command[%d] cannot be empty", buttonKey, i)
			}
		}

		if buttonCfg.Timeout <= 0 {
			return errors.New("config: buttons." + buttonKey + ".timeout must be > 0 (e.g. \"10s\")")
		}
	}

	return nil
}

func isValidLogLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}

func isValidHTTPCodec(codec string) bool {
	switch codec {
	case "event_json", "text", "loki", "ndjson":
		return true
	default:
		return false
	}
}

func isValidHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func validateHostPort(addr string) error {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if portStr == "" {
		return errors.New("port is empty")
	}
	return nil
}

func validateHTTPURL(raw string) error {
	u, err := neturl.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	if u.Host == "" {
		return errors.New("missing host")
	}
	return nil
}
