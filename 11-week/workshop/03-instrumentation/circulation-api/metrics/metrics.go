// Package metrics exposes circulation-api metrics in Prometheus text format.
package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Buckets in seconds: a librarian expects a loan in under 1 s; 5 s is the client timeout.
var Buckets = []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5}

type Metrics struct {
	mu      sync.Mutex
	total   map[[2]string]int  // {route, outcome} -> count
	buckets map[string][]int   // route -> cumulative counts per bucket
	sum     map[string]float64 // route -> sum of seconds
	count   map[string]int
}

func New() *Metrics {
	return &Metrics{total: map[[2]string]int{}, buckets: map[string][]int{},
		sum: map[string]float64{}, count: map[string]int{}}
}

// Observe counts one request by outcome (success, rejected, error) and its latency.
func (m *Metrics) Observe(route, outcome string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.total[[2]string{route, outcome}]++
	if m.buckets[route] == nil {
		m.buckets[route] = make([]int, len(Buckets))
	}
	for i, b := range Buckets {
		if d.Seconds() <= b {
			m.buckets[route][i]++
		}
	}
	m.sum[route] += d.Seconds()
	m.count[route]++
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# TYPE circulation_loan_requests_total counter")
	for k, v := range m.total {
		fmt.Fprintf(w, "circulation_loan_requests_total{route=%q,outcome=%q} %d\n", k[0], k[1], v)
	}
	fmt.Fprintln(w, "# TYPE circulation_loan_request_duration_seconds histogram")
	for route, counts := range m.buckets {
		for i, b := range Buckets {
			fmt.Fprintf(w, "circulation_loan_request_duration_seconds_bucket{route=%q,le=\"%g\"} %d\n", route, b, counts[i])
		}
		fmt.Fprintf(w, "circulation_loan_request_duration_seconds_bucket{route=%q,le=\"+Inf\"} %d\n", route, m.count[route])
		fmt.Fprintf(w, "circulation_loan_request_duration_seconds_sum{route=%q} %g\n", route, m.sum[route])
		fmt.Fprintf(w, "circulation_loan_request_duration_seconds_count{route=%q} %d\n", route, m.count[route])
	}
}
