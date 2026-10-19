package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Hermesss29/w11/circulation-api/catalogclient"
	"github.com/Hermesss29/w11/circulation-api/metrics"
	"github.com/Hermesss29/w11/circulation-api/obslog"
	"github.com/Hermesss29/w11/circulation-api/tracectx"
)

type loanRequest struct {
	StudentID string `json:"studentId"`
	BookID    string `json:"bookId"`
}

// newRouter wires POST /api/v1/loans (HU-06) behind the trace middleware, plus /metrics.
func newRouter(log *slog.Logger, catalog *catalogclient.Client, m *metrics.Metrics) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/loans", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := obslog.FromCtx(r.Context(), log).With("http_route", "/api/v1/loans")
		var req loanRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		l.Info("register loan request received", "event", "loan.registration.started",
			"student_id", req.StudentID, "book_id", req.BookID, "duration_ms", 0)

		status, outcome := http.StatusCreated, "success"
		if err := catalog.LoanCopy(r.Context(), req.BookID); err != nil {
			status, outcome = http.StatusUnprocessableEntity, "rejected"
			l.Warn("loan rejected, no copy reserved", "event", "loan.registration.rejected",
				"error_code", "NO_COPIES_AVAILABLE", "duration_ms", obslog.Since(start))
		}
		w.WriteHeader(status)
		m.Observe("/api/v1/loans", outcome, time.Since(start))
		l.Info("POST /api/v1/loans answered", "event", "http.request.completed",
			"http_status", status, "duration_ms", obslog.Since(start))
	})
	mux.Handle("GET /metrics", m)
	return tracectx.Middleware(mux)
}

func main() {
	log := obslog.New(os.Stdout, "circulation-api", "w11-local", "local")
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	if catalogURL == "" {
		catalogURL = "http://localhost:8082"
	}
	router := newRouter(log, catalogclient.New(catalogURL, log), metrics.New())
	log.Info("circulation-api listening on :8083", "event", "service.started", "duration_ms", 0)
	log.Error("server stopped", "event", "service.stopped", "error", http.ListenAndServe(":8083", router))
}
