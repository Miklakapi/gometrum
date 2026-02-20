package logsinks

import (
	"context"
	"log/slog"
	"time"
)

type SinkAdapter struct {
	minLevel slog.Level
	push     func(LogEvent) bool
	attrs    []slog.Attr
}

func NewSinkAdapter(minLevel slog.Level, push func(LogEvent) bool) *SinkAdapter {
	return &SinkAdapter{
		minLevel: minLevel,
		push:     push,
	}
}

func (a *SinkAdapter) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= a.minLevel
}

func (a *SinkAdapter) Handle(ctx context.Context, r slog.Record) error {
	ev := recordToEvent(r, a.attrs)
	_ = a.push(ev)
	return nil
}

func (a *SinkAdapter) WithAttrs(attrs []slog.Attr) slog.Handler {
	cp := *a
	cp.attrs = append(append([]slog.Attr{}, a.attrs...), attrs...)
	return &cp
}

func (a *SinkAdapter) WithGroup(name string) slog.Handler {
	return a
}

func recordToEvent(r slog.Record, extra []slog.Attr) LogEvent {
	ev := LogEvent{
		Time:  r.Time,
		Level: levelToString(r.Level),
		Msg:   r.Message,
	}

	var attrs map[string]any

	add := func(key string, v any) {
		if attrs == nil {
			attrs = make(map[string]any, 8)
		}
		attrs[key] = v
	}

	var flatten func(prefix string, a slog.Attr)

	flatten = func(prefix string, a slog.Attr) {
		key := a.Key
		if prefix != "" {
			key = prefix + "." + key
		}

		v := a.Value.Resolve()

		if v.Kind() == slog.KindGroup {
			for _, it := range v.Group() {
				flatten(key, it)
			}
			return
		}

		add(key, valueToAnyResolved(v))
	}

	for _, a := range extra {
		flatten("", a)
	}

	r.Attrs(func(a slog.Attr) bool {
		flatten("", a)
		return true
	})

	if attrs != nil {
		ev.Attrs = attrs
	}

	return ev
}

func levelToString(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	default:
		return "error"
	}
}

func valueToAnyResolved(v slog.Value) any {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		return v.Bool()
	case slog.KindInt64:
		return v.Int64()
	case slog.KindUint64:
		return v.Uint64()
	case slog.KindFloat64:
		return v.Float64()
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339Nano)
	case slog.KindAny:
		if err, ok := v.Any().(error); ok {
			return err.Error()
		}
		return v.Any()
	default:
		return v.Any()
	}
}
