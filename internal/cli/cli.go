package cli

import "flag"

type CLI struct {
	ConfigPath     string
	Once           bool
	LogLevel       string
	Quiet          bool
	Validate       bool
	DryRun         bool
	PrintConfig    bool
	GenerateConfig bool
}

func ParseFlags() CLI {
	var cfg CLI

	flag.StringVar(&cfg.ConfigPath, "config", "", "Path to YAML config file (e.g. /etc/gometrum.yaml)")
	flag.StringVar(&cfg.ConfigPath, "c", "", "Shorthand for --config")

	flag.BoolVar(&cfg.Once, "once", false, "Collect and publish once, then exit")

	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level: debug, info, warn, error")

	flag.BoolVar(&cfg.Quiet, "quiet", false, "Suppress all logs except errors")
	flag.BoolVar(&cfg.Quiet, "q", false, "Shorthand for --quiet")

	flag.BoolVar(&cfg.Validate, "validate", false, "Validate config and exit")

	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Collect metrics but do not publish (print output only)")

	flag.BoolVar(&cfg.PrintConfig, "print-config", false, "Print final merged config and exit")

	flag.BoolVar(&cfg.GenerateConfig, "generate-config", false, "Generate example configuration and exit")

	flag.Parse()

	if cfg.Quiet {
		cfg.LogLevel = "error"
	}

	return cfg
}
