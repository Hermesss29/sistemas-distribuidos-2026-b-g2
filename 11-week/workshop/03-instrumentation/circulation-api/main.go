package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Hermesss29/w11/circulation-api/obslog"
	"github.com/Hermesss29/w11/circulation-api/tracectx"
)

type loanRequest struct {
	StudentID string `json:"studentId"`
	BookID    string `json:"bookId"`
}

// newRouter wires POST /api/v1/loans (HU-06) behind the trace middleware.
func newRouter(log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/loans", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		l := obslog.FromCtx(r.Context(), log)
		var req loanRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		l.Info("register loan request received", "event", "loan.registration.started",
			"student_id", req.StudentID, "book_id", req.BookID, "http_route", "/api/v1/loans", "duration_ms", 0)
		w.WriteHeader(http.StatusCreated)
		l.Info("POST /api/v1/loans answered", "event", "http.request.completed",
			"http_route", "/api/v1/loans", "http_status", 201, "duration_ms", obslog.Since(start))
	})
	return tracectx.Middleware(mux)
}

func main() {
	log := obslog.New(os.Stdout, "circulation-api", "w11-local", "local")
	log.Info("circulation-api listening on :8083", "event", "service.started", "duration_ms", 0)
	log.Error("server stopped", "event", "service.stopped", "error", http.ListenAndServe(":8083", newRouter(log)))
}
