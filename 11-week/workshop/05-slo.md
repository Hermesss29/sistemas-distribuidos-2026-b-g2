# CP6 - SLIs, SLOs and error budget

Journey: HU-06 Register a loan (`POST /api/v1/loans` in `circulation-api`). All PromQL uses the
metrics defined in CP4 (`04-metrics.md`) and exposed by the CP3 code.

## Traffic estimate

| Data | Value | Why |
|---|---|---|
| Library counters | 1 counter, 2 librarians | the LMS replaces the paper loan book of one university library |
| Loans per class day | 120 | about 12 loans per hour during 10 open hours |
| Class days per month | 22 | Monday to Friday; weekends are almost empty |
| `POST /api/v1/loans` per 30 days | 120 x 22 = **2,640** | plus a few retries; we round to 2,640 valid requests |
| Window | 30 days = 43,200 minutes | rolling window, same as the course examples |

This is low traffic. One failed request is 0.04 % of the month, so the budget in requests is small and
alerts need a minimum of traffic to avoid firing on a single error (see CP7).

## SLO 1 - Availability

- **SLI:** good requests / valid requests. Good = answered without a server error. `outcome="rejected"`
  (422 `STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`) is good: the system worked and a
  business rule said no. Valid = all requests to `/api/v1/loans`.
- **Objective:** 99.5 % over 30 days.
- **PromQL (SLI over 30 days):**

```promql
1 - (
  sum(increase(circulation_loan_requests_total{route="/api/v1/loans",outcome="error"}[30d]))
  /
  sum(increase(circulation_loan_requests_total{route="/api/v1/loans"}[30d]))
)
```

- **Error budget:** 100 % - 99.5 % = 0.5 %.
  - Requests: 0.5 % x 2,640 = **13 failed loan registrations** per month.
  - Minutes: 0.5 % x 43,200 = **216 minutes** of full outage per month.

## SLO 2 - Latency

- **SLI:** requests answered in 1 second or less / valid requests. 1 s is the bucket chosen in CP4:
  the librarian has the student waiting at the counter.
- **Objective:** 99 % of loan registrations in 1 s or less over 30 days.
- **PromQL:**

```promql
sum(increase(circulation_loan_request_duration_seconds_bucket{route="/api/v1/loans",le="1"}[30d]))
/
sum(increase(circulation_loan_request_duration_seconds_count{route="/api/v1/loans"}[30d]))
```

- **Error budget:** 1 %.
  - Requests: 1 % x 2,640 = **26 slow loan registrations** per month.
  - Minutes: 1 % x 43,200 = **432 minutes** where every request is slow.

## SLO 3 - Loan consistency (product SLO)

A loan touches two databases: catalog reserves a copy in PostgreSQL and circulation saves the loan in
MongoDB. Without a saga (ADR-008), a failure between both leaves a copy reserved with no loan, and the
book "disappears" from the shelf for every other student. This SLO measures that.

- **SLI:** loans saved / copies reserved by `loan-copy`.
- **Objective:** 99.9 % over 30 days.
- **PromQL:**

```promql
sum(increase(circulation_loan_requests_total{route="/api/v1/loans",outcome="success"}[30d]))
/
sum(increase(catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy",status_class="2xx"}[30d]))
```

- **Error budget:** 0.1 %.
  - Requests: 0.1 % x 2,640 = **2.6, so 2 orphan copies** per month.
  - Minutes: 0.1 % x 43,200 = **43 minutes**.
- It is stricter because each failure is a wrong stock number that someone must fix by hand in
  `catalog.book`.

## Error budget policy (team G2)

| Budget consumed | What the team does |
|---|---|
| 50 % | Warning in the team channel. The owner of the service reviews the errors of the week by `trace_id`. New features continue, but every PR that touches `circulation-api`, `membership-api` or `catalog-api` needs a test for the failure path. |
| 75 % | Only low-risk releases. No changes to the HU-06 flow (`LoanRegistrationService`) without a second reviewer. Top action: fix the main cause found in the logs (for example the missing `Idempotency-Key` on `loan-copy`). |
| 100 % | Feature freeze for the services of the journey until the window recovers. Only fixes, rollbacks and reliability work (timeouts, retries, saga compensation). A blameless postmortem is written and reviewed with the teacher in the next session. |
