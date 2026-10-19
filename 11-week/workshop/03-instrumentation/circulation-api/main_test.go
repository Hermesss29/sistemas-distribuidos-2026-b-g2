package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Hermesss29/w11/circulation-api/catalogclient"
	"github.com/Hermesss29/w11/circulation-api/metrics"
	"github.com/Hermesss29/w11/circulation-api/obslog"
	"github.com/Hermesss29/w11/circulation-api/tracectx"
)

// The trace_id that reaches circulation-api must be the same one sent to catalog-api.
func TestTraceIDIsPropagatedToCatalog(t *testing.T) {
	var downstream string
	catalog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstream = r.Header.Get(tracectx.Header)
		w.WriteHeader(http.StatusOK)
	}))
	defer catalog.Close()

	log := obslog.New(io.Discard, "circulation-api", "test", "test")
	m := metrics.New()
	api := httptest.NewServer(newRouter(log, catalogclient.New(catalog.URL, log), m))
	defer api.Close()

	sent := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	req, _ := http.NewRequest(http.MethodPost, api.URL+"/api/v1/loans",
		strings.NewReader(`{"studentId":"5f2c8e1a","bookId":"a41d2f6b"}`))
	req.Header.Set(tracectx.Header, sent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode)
	}

	in, _ := tracectx.Parse(sent)
	out, err := tracectx.Parse(downstream)
	if err != nil {
		t.Fatalf("catalog-api got no valid traceparent: %v", err)
	}
	if out.TraceID != in.TraceID {
		t.Fatalf("trace_id downstream = %s, want %s", out.TraceID, in.TraceID)
	}
	if out.SpanID == in.SpanID {
		t.Fatal("downstream call must use a new span_id")
	}

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, nil)
	if !strings.Contains(rec.Body.String(), `circulation_loan_requests_total{route="/api/v1/loans",outcome="success"} 1`) {
		t.Fatalf("metric not exposed:\n%s", rec.Body.String())
	}
}
