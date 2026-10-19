# CP8 - Runbook: LoanRegistrationAvailabilityFastBurn

| Field | Value |
|---|---|
| Alert | `LoanRegistrationAvailabilityFastBurn` (`06-alerts.yml`), `severity=page`, `team=g2-lms` |
| Journey | HU-06 Register a loan, `POST /api/v1/loans` in `circulation-api` |
| SLO | 99.5 % availability in 30 days; this alert means 2 % of the monthly budget is already spent |
| Environment in this runbook | `lms-qa` local compose (gateway on :8080, circulation-api :8083, membership-api :8081, catalog-api :8082) |

## 1. Symptoms

- Librarians get an error in **Loans -> Register loan**; the request returns 5xx with a `traceId` in the body.
- More than 7.2 % of `POST /api/v1/loans` fail in the last hour and also in the last 5 minutes.
- Students wait at the counter; the book cannot leave the library.

## 2. Confirm (5 minutes)

1. Check the burn is real and still happening:

```promql
circulation:loan_error_ratio:rate5m
sum by (outcome) (rate(circulation_loan_requests_total{route="/api/v1/loans"}[5m]))
```

2. Find which dependency fails, using the downstream latency from CP4:

```promql
histogram_quantile(0.95, sum by (le, downstream) (rate(circulation_downstream_request_duration_seconds_bucket[5m])))
```

3. Take one `traceId` from a failed response (or from the librarian's screen) and read the whole story,
   following the CP2 contract:

```
docker compose logs --since 15m circulation-api membership-api catalog-api | findstr <trace_id>
```

   Look for `event` = `catalog.loan_copy.timeout`, `membership.*.timeout` or an `ERROR` from MongoDB.

4. Check the containers:

```
docker compose ps -a
docker compose logs --tail 50 <service>
```

## 3. Mitigate (in order of risk, lowest first)

Our architecture has no feature flag for HU-06 and no circuit breaker in `circulation-api` yet, so those
options are not available. What we can really do:

| # | Action | Risk | When |
|---|---|---|---|
| 1 | **Restart the failing dependency**: `docker compose restart catalog-api` (or `membership-api`). Then check `/health/ready`. | Low: only that service is down for a few seconds | logs show timeouts or 5xx from one dependency |
| 2 | **Roll back `circulation-api`** to the last good release (tag `v2.0.x` image) with `docker compose up -d circulation-api` on the previous image tag. | Medium: undoes the last deploy; in `main` it needs the teacher's PR approval | errors started right after a deploy of `circulation-api` |
| 3 | **Escalate** to the owner of the failing service and switch the counter to the manual workaround: write the loan on paper and register it later in HU-06. | Low for data, high for effort | 1 and 2 did not stop the burn in 15 minutes |

After mitigating, look for orphan copies (copy reserved without a loan, SLO 3). The fix is manual in
`catalog.book` and must be reviewed by the catalog owner.

## 4. Escalation (team G2)

| Order | Role | What they do |
|---|---|---|
| 1 | On-call member of the week | runs this runbook, writes in the team channel what was done and the `trace_id` |
| 2 | Owner of the failing service (Circulation, Catalog or Membership; Membership: Hermes) | reads the logs, decides restart or rollback, fixes the code |
| 3 | Team lead of G2 | decides the rollback of a release and the error budget policy (`05-slo.md`) |
| 4 | Teacher (approver of PRs to `main`) | approves the hotfix or rollback PR to `main` |

## 5. Verify recovery

- `circulation:loan_error_ratio:rate5m` below 0.5 % for 15 minutes, and the alert is resolved.
- One real test loan from the portal: 201, the book shows one copy less in **Catalog**, and every log
  line of that request has the same `trace_id`.
- Write the timeline (start, detection, mitigation, end) for the postmortem, as in `08-incident.md`.
