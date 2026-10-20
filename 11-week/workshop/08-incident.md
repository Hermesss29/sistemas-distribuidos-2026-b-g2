# CP9 - Incident drill and blameless postmortem

## Scenario

Journey HU-06 Register a loan. The downstream dependency of hop 5 in CP1 (`circulation-api` ->
`catalog-api`, `POST /api/v1/books/{id}/loan-copy`) answers in **3 seconds for 25 % of the requests
during 40 minutes**. Drill run on 2026-10-20, incident from 15:00 to 15:40 (T+0 to T+40).

## Trace of one slow request (span tree with durations)

```
[root] gateway POST /api/v1/loans                                   3165 ms  trace_id=5d0c8a3e7b1f42c69e4a2d7b0c1f8e63
  └─ circulation-api POST /api/v1/loans (server)                    3159 ms
       └─ LoanRegistrationService.RegisterLoan (internal)           3152 ms
            ├─ membershipclient GET /students/{id} (client)           18 ms
            │    └─ membership-api GET /api/v1/students/{id}          11 ms
            ├─ mongo loans.countDocuments (client)                     6 ms
            ├─ catalogclient POST /books/{id}/loan-copy (client)    3118 ms   <- slow span
            │    └─ catalog-api POST /api/v1/books/{id}/loan-copy   3104 ms
            │         ├─ wait for DB connection (pool)              3091 ms
            │         └─ postgres UPDATE catalog.book                  9 ms
            └─ mongo loans.insertOne (client)                          9 ms
```

98 % of the time is in one span: catalog-api waiting for a free database connection, not the query itself.

## Four log lines (CP2 contract)

```json
{"timestamp":"2026-10-20T15:14:07.512Z","level":"INFO","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"5d0c8a3e7b1f42c69e4a2d7b0c1f8e63","span_id":"a1c4e7b29d3f5061","correlation_id":"9a7e2c14-3b6d-4f80-a1e5-6c2d9b8f0e37","event":"loan.registration.started","message":"register loan request received","duration_ms":0,"student_id":"7c3d1e9a-2f4b-4a6c-8e0d-5b1a9f3c7d26","book_id":"e2b9c4d1-6a3f-4e8b-9c7d-1f0a5b2e8c43","http_route":"/api/v1/loans"}
{"timestamp":"2026-10-20T15:14:10.639Z","level":"WARN","service":"catalog-api","version":"eef2ea5","environment":"qa","trace_id":"5d0c8a3e7b1f42c69e4a2d7b0c1f8e63","span_id":"6e2b8d4f1a9c3057","correlation_id":"9a7e2c14-3b6d-4f80-a1e5-6c2d9b8f0e37","event":"book.loan_copy.slow","message":"loan-copy waited 3.1s for a database connection","duration_ms":3104,"book_id":"e2b9c4d1-6a3f-4e8b-9c7d-1f0a5b2e8c43","http_route":"/api/v1/books/{id}/loan-copy","http_status":200}
{"timestamp":"2026-10-20T15:14:10.648Z","level":"WARN","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"5d0c8a3e7b1f42c69e4a2d7b0c1f8e63","span_id":"d7f1a3c95b2e4086","correlation_id":"9a7e2c14-3b6d-4f80-a1e5-6c2d9b8f0e37","event":"catalog.loan_copy.completed","message":"catalog-api answered over the 1 s target","duration_ms":3118,"book_id":"e2b9c4d1-6a3f-4e8b-9c7d-1f0a5b2e8c43","downstream":"catalog-api","http_route":"/api/v1/books/{id}/loan-copy","http_status":200}
{"timestamp":"2026-10-20T15:14:10.671Z","level":"INFO","service":"circulation-api","version":"5586f3f","environment":"qa","trace_id":"5d0c8a3e7b1f42c69e4a2d7b0c1f8e63","span_id":"a1c4e7b29d3f5061","correlation_id":"9a7e2c14-3b6d-4f80-a1e5-6c2d9b8f0e37","event":"http.request.completed","message":"POST /api/v1/loans answered","duration_ms":3159,"http_route":"/api/v1/loans","http_status":201}
```

The request did not fail (201), so only the latency SLO is hurt. The 3 s is below the 5 s timeout of
`catalogclient`, so there is no `DOWNSTREAM_TIMEOUT` and no orphan copy.

## SLIs during the incident

Traffic from CP6: about 12 loans per hour, so 40 minutes = **8 loan registrations**, 25 % = **2 slow**.

| SLO | Value during the 40 min | Value over 30 days | Budget consumed |
|---|---|---|---|
| Availability 99.5 % | 100 % (no 5xx) | not affected | 0 % |
| Latency 99 % under 1 s | 75 % (25 % slow) | 1 - 2/2,640 = 99.92 % | 2 of 26 slow requests = **7.7 %** by requests; 40 min x 25 % = 10 min of 432 = **2.3 %** by time |
| Loan consistency 99.9 % | 100 % | not affected | 0 % |

## Which alert fires first, and when

- Only `LoanRegistrationLatencyFastBurn` can fire. The availability alerts never fire, because there are
  no errors.
- Burn rate of latency = 25 % / 1 % = **25**.
- The 5 m window passes 14.4 % at about T+5 (25 % slow).
- The 1 h window averages in the 59 good minutes before the incident:
  - slow ratio(1h) = 25 % x t / 60;
  - it passes 14.4 % when t > 34.6 minutes, so at **T+35**.
- With `for: 2m`, the page arrives at **T+37 (15:37)**.
- The slow burn (6x over 6h) needs 25 % x t / 360 > 6 % -> t > 86 min, so it never fires in 40 minutes.

## Runbook step that stops the damage

- T+38: step 2 of the runbook (downstream latency) shows `catalog-api` p95 = 3.1 s, while `membership-api`
  is 0.02 s.
- T+39: step 3 (logs by `trace_id`) shows `book.loan_copy.slow` "waited for a database connection".
- T+40: **mitigation 1, `docker compose restart catalog-api`**, frees the pool. The next loans take 0.2 s.
- T+47: the alert resolves.

## Postmortem (blameless)

**Timeline**

| Time | T+ | Event |
|---|---|---|
| 15:00 | 0 | a long list export in the catalog portal holds most catalog-api DB connections; 1 of 4 `loan-copy` calls waits about 3 s |
| 15:05 | 5 | 5 m slow ratio above 14.4 %; nobody notices, the loans still succeed |
| 15:37 | 37 | `LoanRegistrationLatencyFastBurn` pages the on-call member |
| 15:38 | 38 | runbook step 2 points to catalog-api |
| 15:39 | 39 | logs by `trace_id` show the wait for a database connection |
| 15:40 | 40 | `catalog-api` restarted; latency back to normal (end of impact) |
| 15:47 | 47 | alert resolved; test loan 201 in 0.2 s, same `trace_id` in all lines |

**Impact:**

- 2 of 8 loan registrations took about 3 s, with the student waiting at the counter.
- No failed loans, and no orphan copies.
- 7.7 % of the monthly latency budget by requests (2.3 % by time).

**Root cause:** in `catalog-api` the HU-05 book search and `loan-copy` share one small connection pool.
A long search held the connections, and `loan-copy` queued behind it.

**Contributing factors:**

- The PostgreSQL instance is shared by the domains (Annex J).
- The `catalogclient` timeout is 5 s, longer than the 1 s target, so the slow call never failed fast.
- The real `circulation-api` has no downstream latency metric yet (only this workshop code).
- With about 12 loans per hour, the 1 h window needs 35 minutes to reach the threshold.

**Actions:**

1. **As a librarian, I want a loan to be registered in under 1 s even when someone is searching the
   catalog, so that the student does not wait at the counter.**
   - AC1: `catalog-api` uses a separate connection pool, or a reserved quota, for `loan-copy`.
   - AC2: a load test with 20 concurrent searches keeps `loan-copy` p95 under 200 ms.
2. **As the on-call member, I want `circulation-api` to expose `circulation_downstream_request_duration_seconds`
   by `downstream`, so that I find the slow dependency in one query.**
   - AC1: `/metrics` of the real `circulation-api` shows the histogram for `catalog-api` and `membership-api`.
   - AC2: the runbook query in step 2 returns data in `lms-qa`.
3. **As the circulation owner, I want `catalogclient` to fail fast and retry once with an `Idempotency-Key`,
   so that a slow catalog does not block the counter and no copy is reserved twice.**
   - AC1: the timeout is 1.5 s, with one retry and the same `Idempotency-Key`.
   - AC2: a test with a 3 s fake catalog gets an answer from `circulation-api` in under 3.5 s, and the
     fake receives the same key twice.
