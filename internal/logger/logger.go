package logger

import (
	"log/slog"
	"os"
	"strings"
)

func SetupLogger(level string) {
	normalized := strings.ToLower(level)

	var lvl slog.Level
	switch normalized {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelWarn
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	slog.SetDefault(slog.New(handler))

	if normalized != "debug" && normalized != "info" && normalized != "warn" && normalized != "error" {
		slog.Warn("Unknown log level, falling back to warn", "provided", level)
	}
}
