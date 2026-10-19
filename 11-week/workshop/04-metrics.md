# CP4 - Golden signals and metric catalog

Journey: HU-06 Register a loan. One table per service the request crosses. Names follow Prometheus
conventions: snake_case, service prefix, unit suffix (`_seconds`), `_total` for counters.
`circulation_loan_requests_total` and `circulation_loan_request_duration_seconds` are already exposed
by the CP3 code on `/metrics`.

## circulation-api (orchestrator of HU-06)

| Signal | Metric | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `circulation_loan_request_duration_seconds` | histogram | `route` | `histogram_quantile(0.95, sum by (le, route) (rate(circulation_loan_request_duration_seconds_bucket[5m])))` |
| Traffic | `circulation_loan_requests_total` | counter | `route`, `outcome` | `sum by (route) (rate(circulation_loan_requests_total[5m]))` |
| Errors | `circulation_loan_requests_total` | counter | `route`, `outcome` | `sum(rate(circulation_loan_requests_total{outcome="error"}[5m])) / sum(rate(circulation_loan_requests_total[5m]))` |
| Saturation | `circulation_http_inflight_requests` | gauge | none | `max_over_time(circulation_http_inflight_requests[5m])` |
| Saturation (dependency) | `circulation_downstream_request_duration_seconds` | histogram | `downstream`, `outcome` | `histogram_quantile(0.95, sum by (le, downstream) (rate(circulation_downstream_request_duration_seconds_bucket[5m])))` |

`outcome` values: `success` (201), `rejected` (422 business rule: `STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`,
`NO_COPIES_AVAILABLE`), `error` (5xx or downstream timeout). A 422 is not an error for the SLO: the
system worked and a rule said no.

## membership-api (hop 3: `GET /api/v1/students/{id}`)

| Signal | Metric | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `membership_http_request_duration_seconds` | histogram | `route`, `method` | `histogram_quantile(0.95, sum by (le, route) (rate(membership_http_request_duration_seconds_bucket[5m])))` |
| Traffic | `membership_http_requests_total` | counter | `route`, `method`, `status_class` | `sum by (route) (rate(membership_http_requests_total[5m]))` |
| Errors | `membership_http_requests_total` | counter | `route`, `method`, `status_class` | `sum(rate(membership_http_requests_total{status_class="5xx"}[5m]))` |
| Saturation | `membership_db_connections_in_use` | gauge | none | `max_over_time(membership_db_connections_in_use[5m]) / membership_db_connections_max` |

## catalog-api (hop 5: `POST /api/v1/books/{id}/loan-copy`)

| Signal | Metric | Type | Labels | PromQL (dashboard) |
|---|---|---|---|---|
| Latency | `catalog_http_request_duration_seconds` | histogram | `route`, `method` | `histogram_quantile(0.95, sum by (le) (rate(catalog_http_request_duration_seconds_bucket{route="/api/v1/books/{id}/loan-copy"}[5m])))` |
| Traffic | `catalog_http_requests_total` | counter | `route`, `method`, `status_class` | `sum by (route) (rate(catalog_http_requests_total[5m]))` |
| Errors | `catalog_http_requests_total` | counter | `route`, `method`, `status_class` | `sum(rate(catalog_http_requests_total{status_class="5xx"}[5m]))` |
| Saturation | `catalog_db_connections_in_use` | gauge | none | `max_over_time(catalog_db_connections_in_use[5m]) / catalog_db_connections_max` |
| Business | `catalog_loan_copy_rejected_total` | counter | `reason` | `sum by (reason) (increase(catalog_loan_copy_rejected_total[1h]))` |

## Label cardinality

| Label | Possible values | Estimate |
|---|---|---|
| `route` | route templates of each API (`/api/v1/loans`, `/api/v1/loans/{id}/return`, `/api/v1/students/{id}`, ...) | about 6 per service |
| `method` | GET, POST, PATCH | 3 |
| `outcome` | success, rejected, error | 3 |
| `status_class` | 2xx, 4xx, 5xx | 3 |
| `downstream` | membership-api, catalog-api, mongodb | 3 |
| `reason` | `NO_COPIES_AVAILABLE`, `NOT_FOUND` | 2 |
| `le` | 8 buckets including `+Inf` | 8 |

Worst series count: `circulation_loan_request_duration_seconds` = 6 routes x 8 buckets = 48 series;
`membership_http_requests_total` = 6 x 3 x 3 = 54 series. Small and bounded.

**Forbidden labels** (unbounded, they would create one series per student or book):

| Forbidden | Why | Use instead |
|---|---|---|
| `student_id`, `document_id`, `email` | one series per student, and personal data | logs with `trace_id` (CP2) |
| `book_id`, `loan_id` | one series per book or loan | logs and the `reason` label |
| raw URL (`/api/v1/students/5f2c...`) | ids inside the path | `route` template (`/api/v1/students/{id}`) |
| `trace_id`, `correlation_id`, `Idempotency-Key` | unique per request | exemplars or logs |

## Latency buckets

`0.05, 0.1, 0.25, 0.5, 1, 2.5, 5` seconds (`metrics.Buckets` in CP3).

- The user is a librarian at the counter with the student waiting. A loan registered in under 1 s
  feels instant; between 1 s and 2.5 s it is acceptable; more than that the librarian clicks again.
- HU-06 makes 2 HTTP calls and 2 MongoDB operations. Each one is usually under 50 ms in local, so
  `0.05` to `0.25` separate the normal cases.
- `5` is the `http.Client` timeout of `membershipclient` and `catalogclient`. Anything that reaches it
  is a timeout, so there is no reason for bigger buckets.
- `1` is the bucket the latency SLO will use in CP6 (for example: 95 % of loans under 1 s).
