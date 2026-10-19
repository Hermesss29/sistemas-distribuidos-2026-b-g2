// Package catalogclient calls catalog-api and propagates the trace context (hop 5 of HU-06).
package catalogclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Hermesss29/w11/circulation-api/obslog"
	"github.com/Hermesss29/w11/circulation-api/tracectx"
)

type Client struct {
	baseURL string
	http    *http.Client
	log     *slog.Logger
}

func New(baseURL string, log *slog.Logger) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 5 * time.Second}, log: log}
}

// LoanCopy reserves one copy with POST /api/v1/books/{id}/loan-copy.
func (c *Client) LoanCopy(ctx context.Context, bookID string) error {
	start := time.Now()
	child := tracectx.From(ctx).Child()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/api/v1/books/%s/loan-copy", c.baseURL, bookID), nil)
	if err != nil {
		return err
	}
	req.Header.Set(tracectx.Header, child.Traceparent())
	req.Header.Set(tracectx.CorrelationHeader, child.CorrelationID)

	l := obslog.FromCtx(tracectx.With(ctx, child), c.log).With("downstream", "catalog-api",
		"book_id", bookID, "http_route", "/api/v1/books/{id}/loan-copy")
	resp, err := c.http.Do(req)
	if err != nil {
		l.Warn("catalog-api did not answer", "event", "catalog.loan_copy.timeout",
			"error_code", "DOWNSTREAM_TIMEOUT", "duration_ms", obslog.Since(start))
		return err
	}
	defer resp.Body.Close()
	l.Info("catalog-api answered", "event", "catalog.loan_copy.completed",
		"http_status", resp.StatusCode, "duration_ms", obslog.Since(start))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("catalog-api returned %d", resp.StatusCode)
	}
	return nil
}
