# How to run - circulation-api instrumentation (CP3)

Minimal module of `circulation-api` for the HU-06 loan journey. It only uses the Go standard library
(`net/http`, `log/slog`), so there is nothing to install besides Go 1.22 or newer.

| Folder | What it does |
|---|---|
| `tracectx/` | reads or creates the W3C `traceparent`, opens a child span, returns it in the response |
| `obslog/` | JSON logger with the CP2 contract fields (`trace_id`, `span_id`, `correlation_id`, ...) |
| `catalogclient/` | `POST /api/v1/books/{id}/loan-copy` to catalog-api, copying `traceparent` and `X-Correlation-Id` |
| `metrics/` | `circulation_loan_requests_total` and `circulation_loan_request_duration_seconds` on `/metrics` |
| `main_test.go` | proves the `trace_id` received is the same one sent to catalog-api |

## Install

```
cd 11-week/workshop/03-instrumentation/circulation-api
go version
go vet ./...
```

## Test

```
go test ./... -v
```

Real output on my machine (Windows 11, 2026-10-19):

```
=== RUN   TestTraceIDIsPropagatedToCatalog
--- PASS: TestTraceIDIsPropagatedToCatalog (0.00s)
PASS
ok      github.com/Hermesss29/w11/circulation-api       1.810s
?       github.com/Hermesss29/w11/circulation-api/catalogclient [no test files]
?       github.com/Hermesss29/w11/circulation-api/metrics       [no test files]
?       github.com/Hermesss29/w11/circulation-api/obslog        [no test files]
?       github.com/Hermesss29/w11/circulation-api/tracectx      [no test files]
```

## Run

`CATALOG_SERVICE_URL` defaults to `http://localhost:8082` (catalog-api in the lms-qa compose).

```
set CATALOG_SERVICE_URL=http://localhost:8082
go run .
```

Real output on my machine:

```
{"timestamp":"2026-10-19T14:35:38.686Z","level":"INFO","message":"circulation-api listening on :8083","service":"circulation-api","version":"w11-local","environment":"local","event":"service.started","duration_ms":0}
```

Note: when stopping with Ctrl+C on Windows, Go prints `exit status 0xc000013a`. It is the interrupt
signal, not an error.

## Try it by hand

In another terminal (stop the lms-qa `circulation-api` container first, it uses the same port 8083):

```
curl -i -X POST http://localhost:8083/api/v1/loans -H "traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01" -d "{\"studentId\":\"s1\",\"bookId\":\"b1\"}"
curl http://localhost:8083/metrics
```

The response header `traceparent` keeps trace_id `4bf92f3577b34da6a3ce929d0e0e4736` with a new span_id,
and every log line of the request has that same `trace_id`.
