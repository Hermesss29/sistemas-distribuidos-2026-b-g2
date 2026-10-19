// Package obslog writes the CP2 log contract (log/slog JSON).
package obslog

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/Hermesss29/w11/circulation-api/tracectx"
)

func New(out io.Writer, service, version, env string) *slog.Logger {
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String("timestamp", a.Value.Time().UTC().Format("2006-01-02T15:04:05.000Z"))
			case slog.MessageKey:
				a.Key = "message"
			}
			return a
		},
	})
	return slog.New(h).With("service", service, "version", version, "environment", env)
}

// FromCtx injects trace_id, span_id and correlation_id of the current request.
func FromCtx(ctx context.Context, l *slog.Logger) *slog.Logger {
	c := tracectx.From(ctx)
	return l.With("trace_id", c.TraceID, "span_id", c.SpanID, "correlation_id", c.CorrelationID)
}

func Since(start time.Time) int64 { return time.Since(start).Milliseconds() }
