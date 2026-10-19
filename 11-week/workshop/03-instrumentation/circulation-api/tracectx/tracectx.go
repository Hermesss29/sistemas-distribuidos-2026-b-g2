// Package tracectx reads or creates the W3C trace context for circulation-api.
package tracectx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const (
	Header            = "traceparent"
	CorrelationHeader = "X-Correlation-Id"
)

// Context is one span: shared trace_id plus the id of this hop.
type Context struct {
	TraceID       string
	SpanID        string
	CorrelationID string
}

type key struct{}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Parse reads "00-<32 hex trace_id>-<16 hex parent span>-<flags>".
func Parse(h string) (Context, error) {
	p := strings.Split(h, "-")
	if len(p) != 4 || p[0] != "00" || len(p[1]) != 32 || len(p[2]) != 16 {
		return Context{}, fmt.Errorf("invalid traceparent %q", h)
	}
	return Context{TraceID: p[1], SpanID: p[2]}, nil
}

// New starts a trace when no traceparent arrives (first hop).
func New() Context { return Context{TraceID: randomHex(16), SpanID: randomHex(8)} }

// Child keeps the trace_id and opens a new span for the next hop.
func (c Context) Child() Context {
	return Context{TraceID: c.TraceID, SpanID: randomHex(8), CorrelationID: c.CorrelationID}
}

func (c Context) Traceparent() string { return "00-" + c.TraceID + "-" + c.SpanID + "-01" }

func With(ctx context.Context, c Context) context.Context { return context.WithValue(ctx, key{}, c) }
func From(ctx context.Context) Context                    { c, _ := ctx.Value(key{}).(Context); return c }

// Middleware reads traceparent (or creates one), opens the server span of
// circulation-api and returns the context in the response headers.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := Parse(r.Header.Get(Header))
		if err != nil {
			c = New()
		} else {
			c = c.Child()
		}
		c.CorrelationID = r.Header.Get(CorrelationHeader)
		if c.CorrelationID == "" {
			c.CorrelationID = c.TraceID
		}
		w.Header().Set(Header, c.Traceparent())
		w.Header().Set(CorrelationHeader, c.CorrelationID)
		next.ServeHTTP(w, r.WithContext(With(r.Context(), c)))
	})
}
