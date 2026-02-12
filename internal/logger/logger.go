package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Miklakapi/gometrum/internal/config"
)

func SetupBootstrap(level string) {
	lvl, normalized, ok := parseLevel(level)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	slog.SetDefault(slog.New(handler))

	if !ok {
		slog.Warn("Unknown log level, falling back to info", "provided", level, "normalized", normalized)
	}
}

func SetupLogger(cfg config.LogConfig) {
	lvl, normalized, ok := parseLevel(cfg.Level)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	slog.SetDefault(slog.New(handler))

	if !ok {
		slog.Warn("Unknown log level in config, falling back to info", "provided", cfg.Level, "normalized", normalized)
	}

	// sinks (UDP/HTTP) will be attached here in next step.
}

func parseLevel(level string) (slog.Level, string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(level))

	switch normalized {
	case "debug":
		return slog.LevelDebug, normalized, true
	case "info", "":
		return slog.LevelInfo, normalized, normalized != ""
	case "warn", "warning":
		return slog.LevelWarn, normalized, true
	case "error":
		return slog.LevelError, normalized, true
	default:
		return slog.LevelInfo, normalized, false
	}
}
