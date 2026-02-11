package cli

import (
	"errors"
	"flag"
	"strings"
)

type CLI struct {
	ConfigPath      string
	Once            bool
	DryRun          bool
	Validate        bool
	ServicePath     string
	GenerateConfig  bool
	GenerateService bool
	Purge           bool
	Version         bool
}

func ParseFlags() (CLI, error) {
	var cfg CLI

	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to YAML config file (e.g. /etc/gometrum.yaml)")
	flag.StringVar(&cfg.ConfigPath, "c", "", "Shorthand for --config")

	flag.BoolVar(&cfg.Once, "once", false, "Collect and publish once, then exit")

	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Collect metrics but do not publish (print output only)")

	flag.BoolVar(&cfg.Validate, "validate", false, "Validate config and exit")

	flag.BoolVar(&cfg.GenerateConfig, "generate-config", false, "Generate example configuration and exit")

	flag.StringVar(&cfg.ServicePath, "service-path", "", "Path for generated systemd service file (if empty, prints to stdout)")

	flag.BoolVar(&cfg.GenerateService, "generate-service", false, "Generate example systemd service file and exit (prints to stdout if --service-path is empty)")

	flag.BoolVar(&cfg.Purge, "purge", false, "Purge Home Assistant MQTT discovery entities defined in config (publish empty retained configs) and exit")

	flag.BoolVar(&cfg.Version, "version", false, "Show version and exit")
	flag.BoolVar(&cfg.Version, "v", false, "Shorthand for --version")

	flag.Parse()

	if strings.TrimSpace(cfg.ConfigPath) == "" && !cfg.GenerateConfig {
		cfg.ConfigPath = "./gometrum.yaml"
	}

	if err := validateFlags(cfg); err != nil {
		return CLI{}, err
	}

	return cfg, nil
}

func validateFlags(c CLI) error {
	exitModes := 0
	if c.GenerateConfig {
		exitModes++
	}
	if c.Validate {
		exitModes++
	}
	if c.Purge {
		exitModes++
	}
	if c.GenerateService {
		exitModes++
	}
	if c.Version {
		exitModes++
	}

	if exitModes > 1 {
		return errors.New("choose only one of: --generate-config, --validate, --purge, --generate-service, --version,")
	}

	if exitModes > 0 && (c.Once || c.DryRun) {
		return errors.New("flags --once and --dry-run cannot be used with --generate-config, --validate, --purge, --generate-service, --version,")
	}

	return nil
}
