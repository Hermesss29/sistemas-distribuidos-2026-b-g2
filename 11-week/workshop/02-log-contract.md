# CP2 - Structured log contract

All LMS Go services (`access-api`, `membership-api`, `catalog-api`, `circulation-api`) already log JSON
with zap (`internal/infrastructure/logger`). This contract fixes the field names so one query by
`trace_id` returns the whole HU-06 loan registration across services.

One JSON object per line, UTF-8, written to stdout. Field names in snake_case.

## Required fields

| Field | Type | Required | Example | Why it exists |
|---|---|---|---|---|
| `timestamp` | string, ISO-8601 UTC, ms | yes | `2026-10-19T14:20:31.104Z` | order events across containers; zap default `ts` epoch is renamed |
| `level` | string | yes | `INFO` | filter by severity (`DEBUG`, `INFO`, `WARN`, `ERROR`) |
| `service` | string | yes | `circulation-api` | which container wrote the line |
| `version` | string | yes | `5586f3f` | git tag or short commit, to link an error to a release |
| `environment` | string | yes | `qa` | `local`, `qa`, `prod`; avoids mixing QA noise with prod |
| `trace_id` | string, 32 hex | yes | `7c1e4b2a9f3d48e6a0b5c2d1e8f97a34` | joins all hops of one request (from `traceparent`) |
| `span_id` | string, 16 hex | yes | `b3a1c9d2e4f60718` | which hop inside the trace wrote the line |
| `correlation_id` | string, UUID | yes | `d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22` | value of `X-Correlation-Id`; it is the `traceId` the user sees in error bodies |
| `event` | string, dot.case | yes | `loan.registration.rejected` | stable name to count and alert on, never free text |
| `message` | string | yes | `student reached the active loan limit` | short human description |
| `duration_ms` | integer | yes on `*.completed`, `*.timeout`, outbound calls | `5003` | latency per hop without a tracing backend |

## Domain fields (LMS)

| Field | Type | Required | Example | Why it exists |
|---|---|---|---|---|
| `loan_id` | string, UUID | when known | `0b9e...` | follow one loan from HU-06 to its return (HU-07/HU-08) |
| `student_id` | string, UUID | on loan and student events | `5f2c...` | internal id, not personal data by itself; never the document number |
| `book_id` | string, UUID | on loan and catalog events | `a41d...` | find which title has no copies or fails on `loan-copy` |
| `active_loans` | integer | on `loan.active_count.checked` | `2` | explains INV-004 rejections (max 2) |
| `error_code` | string | on WARN/ERROR | `LOAN_LIMIT_REACHED` | same code the API returns (`STUDENT_SUSPENDED`, `NO_COPIES_AVAILABLE`, ...) |
| `http_route` | string | on HTTP events | `/api/v1/books/{id}/loan-copy` | route template, never the raw URL with ids |
| `http_status` | integer | on HTTP events | `422` | result of the hop |
| `downstream` | string | on outbound calls | `catalog-api` | which dependency was called |

## Never log (and how to mask)

| Data | Where it appears in LMS | Rule |
|---|---|---|
| Student full name, email, phone | `membership.student`, HU-02 body | do not log; use `student_id` |
| Student document number (`document_id`) | HU-02, HU-06 search | do not log; if needed for support, last 3 digits: `***456` |
| Access token (RS256 JWT) | `Authorization: Bearer` | never log headers; log only `auth.result=ok/expired/invalid` |
| Internal service token (HS256, `INTERNAL_JWT_SECRET`) | membershipclient, catalogclient | never log; not even the length |
| Admin password / bcrypt hash | `access-api`, `administrators` table | never log the login body |
| DB connection strings | `DB_PASSWORD`, Mongo URI | log host and db name only, without user and password |
| `Idempotency-Key` | POST headers | log only the first 8 chars to debug retries |

## Example: one failed loan registration (same `trace_id` in the seven lines)

Scenario: `catalog-api` is slow. `catalogclient` (with one retry, proposed in CP3) times out on the first
`loan-copy`, retries, and the retry is rejected because the first call had already taken the last copy.

```json
{"timestamp":"2026-10-19T14:20:31.104Z","level":"INFO","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"b3a1c9d2e4f60718","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"loan.registration.started","message":"register loan request received","duration_ms":0,"student_id":"5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97","book_id":"a41d2f6b-9c8e-4e1a-b7d3-0f5c2a8e6b14","http_route":"/api/v1/loans"}
{"timestamp":"2026-10-19T14:20:31.139Z","level":"INFO","service":"membership-api","version":"a729d66","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"4e8f2a6c1b9d3075","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"student.fetched","message":"student found, no active suspension","duration_ms":12,"student_id":"5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97","http_route":"/api/v1/students/{id}","http_status":200}
{"timestamp":"2026-10-19T14:20:31.151Z","level":"INFO","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"b3a1c9d2e4f60718","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"loan.active_count.checked","message":"student has 1 active loan","duration_ms":6,"student_id":"5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97","active_loans":1}
{"timestamp":"2026-10-19T14:20:36.158Z","level":"WARN","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"91c7e3b5a2d84f60","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"catalog.loan_copy.timeout","message":"catalog-api did not answer in 5s, retrying once","duration_ms":5003,"book_id":"a41d2f6b-9c8e-4e1a-b7d3-0f5c2a8e6b14","downstream":"catalog-api","http_route":"/api/v1/books/{id}/loan-copy","error_code":"DOWNSTREAM_TIMEOUT"}
{"timestamp":"2026-10-19T14:20:36.912Z","level":"INFO","service":"catalog-api","version":"eef2ea5","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"2d6a9f1e7c3b0845","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"book.loan_copy.rejected","message":"no available copies left","duration_ms":9,"book_id":"a41d2f6b-9c8e-4e1a-b7d3-0f5c2a8e6b14","http_route":"/api/v1/books/{id}/loan-copy","http_status":422,"error_code":"NO_COPIES_AVAILABLE"}
{"timestamp":"2026-10-19T14:20:36.915Z","level":"WARN","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"b3a1c9d2e4f60718","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"loan.registration.rejected","message":"loan rejected by business rule, book has no copies","duration_ms":5811,"student_id":"5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97","book_id":"a41d2f6b-9c8e-4e1a-b7d3-0f5c2a8e6b14","error_code":"NO_COPIES_AVAILABLE"}
{"timestamp":"2026-10-19T14:20:36.917Z","level":"WARN","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"7c1e4b2a9f3d48e6a0b5c2d1e8f97a34","span_id":"b3a1c9d2e4f60718","correlation_id":"d6f1a7a0-2c43-4f0e-9a51-3e8b7c0d1f22","event":"http.request.completed","message":"POST /api/v1/loans answered 422 after a downstream timeout; possible copy reserved without loan","duration_ms":5813,"http_route":"/api/v1/loans","http_status":422,"error_code":"NO_COPIES_AVAILABLE"}
```

Reading the seven lines by `trace_id` shows the real problem: the first `loan-copy` probably succeeded
in `catalog-api` after the timeout, so a copy is reserved with no loan. This is the missing
compensation (saga, ADR-008), and it is why `loan-copy` also needs an `Idempotency-Key`.

## When we use each level

| Level | Rule | LMS example |
|---|---|---|
| `ERROR` | the request failed for a reason that is not the user's fault and someone must look at it | `circulation-api` cannot reach MongoDB `loan_db`; 500 on `POST /api/v1/loans` |
| `WARN` | expected failure or degraded path; the system handled it | 422 `STUDENT_SUSPENDED` / `LOAN_LIMIT_REACHED`, a downstream timeout with retry |
| `INFO` | business event that changes state, one per step of the journey | `loan.registered`, `student.suspended` (HU-08), `student.deactivated` |
| `DEBUG` | detail only for local troubleshooting, off in `qa` and `prod` | the computed `dueDate`, the response body from `membership-api` |
