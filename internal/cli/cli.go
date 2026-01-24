package cli

import (
	"errors"
	"flag"
	"strings"
)

type CLI struct {
	ConfigPath     string
	Once           bool
	LogLevel       string
	Quiet          bool
	DryRun         bool
	Validate       bool
	PrintConfig    bool
	GenerateConfig bool
}

func ParseFlags() (CLI, error) {
	var cfg CLI

	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to YAML config file (e.g. /etc/gometrum.yaml)")
	flag.StringVar(&cfg.ConfigPath, "c", "", "Shorthand for --config")

	flag.BoolVar(&cfg.Once, "once", false, "Collect and publish once, then exit")

	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level: debug, info, warn, error")

	flag.BoolVar(&cfg.Quiet, "quiet", false, "Suppress all logs except errors")
	flag.BoolVar(&cfg.Quiet, "q", false, "Shorthand for --quiet")

	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Collect metrics but do not publish (print output only)")

	flag.BoolVar(&cfg.Validate, "validate", false, "Validate config and exit")

	flag.BoolVar(&cfg.PrintConfig, "print-config", false, "Print final merged config and exit")

	flag.BoolVar(&cfg.GenerateConfig, "generate-config", false, "Generate example configuration and exit")

	flag.Parse()

	if strings.TrimSpace(cfg.ConfigPath) == "" && !cfg.GenerateConfig {
		cfg.ConfigPath = "./gometrum.yaml"
	}

	if cfg.Quiet {
		cfg.LogLevel = "error"
	} else if strings.TrimSpace(cfg.LogLevel) == "" {
		cfg.LogLevel = "warn"
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
	if c.PrintConfig {
		exitModes++
	}
	if c.Validate {
		exitModes++
	}

	if exitModes > 1 {
		return errors.New("choose only one of: --generate-config, --print-config, --validate")
	}

	if exitModes > 0 && (c.Once || c.DryRun) {
		return errors.New("flags --once and --dry-run cannot be used with --generate-config/--print-config/--validate")
	}

	return nil
}
