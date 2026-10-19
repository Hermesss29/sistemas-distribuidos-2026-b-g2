# CP7 - Alert policy

Rules: `06-alerts.yml` (validated with `promtool check rules 06-alerts.yml` -> `SUCCESS: 9 rules found`).
SLOs and budgets come from `05-slo.md`.

## Burn rate in one line

Burn rate = observed error ratio / allowed error ratio. A burn rate of 1 spends exactly the whole budget
in 30 days (720 hours). Budget spent by an alert before it fires = burn rate x long window / 720 h.

## What each alert means with our numbers

| Alert | Burn rate | Windows | Budget spent before firing | In our units | Severity |
|---|---|---|---|---|---|
| `LoanRegistrationAvailabilityFastBurn` | 14.4 | 1h, confirmed 5m | 14.4 x 1 / 720 = **2 %** | 2 % of 216 min = **4.3 min** of full outage; error ratio > 7.2 % | `page` |
| `LoanRegistrationAvailabilitySlowBurn` | 6 | 6h, confirmed 30m | 6 x 6 / 720 = **5 %** | 5 % of 216 min = **10.8 min**; error ratio > 3 % | `ticket` |
| `LoanRegistrationLatencyFastBurn` | 14.4 | 1h, confirmed 5m | **2 %** of the latency budget | 2 % of 432 min = **8.6 min**; more than 14.4 % of loans over 1 s | `page` |

- **Fast burn (page):** at this speed the whole monthly budget is gone in 720 / 14.4 = 50 hours, about
  two days of classes. Someone of the team looks at it the same hour.
- **Slow burn (ticket):** 720 / 6 = 120 hours, about five days. Nobody needs to stop class; it goes to
  the board and is reviewed in the next daily.
- **Short window (5m / 30m):** stops the alert as soon as the problem is fixed, so it does not keep
  firing for an hour after a rollback.

## Low traffic guard

With 2,640 loans per month (about 12 per hour), the budget in requests is only 13 failed loans. In a
quiet hour, 1 error out of 5 requests is already a 20 % error ratio and would page.

- That is why both availability alerts require a minimum number of requests in the long window:
  `>= 5` in 1h for the fast burn, and `>= 10` in 6h for the slow burn.
- `for: 2m` and `for: 15m` avoid firing on one scrape.
- A single isolated failed loan is still visible in the CP2 logs by `trace_id` and in the dashboard. It
  just does not wake anyone up.

## Cause-based alerts we will NOT create

| Alert we skip | Why not |
|---|---|
| CPU of the `circulation-api` container above 80 % | High CPU is not a user problem by itself. Go services can run hot and still answer in 50 ms; if CPU really hurts, latency rises and `LoanRegistrationLatencyFastBurn` fires. A CPU alert would page people for nothing during load tests. |
| `catalog-api` or `membership-api` container restarted / health check failed once | Docker restarts the container and `circulation-api` only fails while it is down. If that hurts loans, the availability burn alerts fire. Alerting on every restart (for example a redeploy of `catalog-api`) creates noise and teaches the team to ignore pages. |

Rule we follow: page on symptoms that the librarian feels (failed or slow loans). Causes (CPU, restarts,
MongoDB connections) go to dashboards and are used during the investigation in the runbook.
