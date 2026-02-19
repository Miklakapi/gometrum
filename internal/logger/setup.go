package logger

import (
	"log/slog"
	"os"
)

func SetupBootstrap(level string) {
	lvl, ok := ParseLevel(level)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	slog.SetDefault(slog.New(handler))

	if !ok {
		slog.Warn("Unknown log level, falling back to info", "provided", level)
	}
}

func Setup(level string, extraHandlers ...slog.Handler) {
	lvl, ok := ParseLevel(level)

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	handlers := make([]slog.Handler, 0, 1+len(extraHandlers))
	handlers = append(handlers, consoleHandler)

	for _, h := range extraHandlers {
		if h != nil {
			handlers = append(handlers, h)
		}
	}

	slog.SetDefault(slog.New(NewMultiHandler(handlers...)))

	if !ok {
		slog.Warn("Unknown log level, falling back to info", "provided", level)
	}
}
func ParseLevel(level string) (slog.Level, bool) {
	switch level {
	case "debug":
		return slog.LevelDebug, true
	case "info", "":
		return slog.LevelInfo, level != ""
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}
