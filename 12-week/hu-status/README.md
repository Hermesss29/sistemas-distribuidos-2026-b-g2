<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       12-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 12

> Note: following the course calendar, week 11 was the project presentation and week 12 was the W11 guided workshop (Monday 2026-10-19 in class, Tuesday 2026-10-20 at home). The workshop guide required the files under `11-week/workshop/` and the reflection in `11-week/hu-status/`, so they stay there with their original push times; this report links to that same work.

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Make the HU-06 loan registration journey of the LMS observable end to end (trace map, structured log contract, trace-context instrumentation for circulation-api with a test, and a golden-signals metric catalog) in the week 11 guided workshop, then define SLOs, burn-rate alerts, a runbook and an incident drill for it (deliveries 1 and 2).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-06 | W11 CP0 - start observability workshop | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/53bd527 |
| HU-06 | W11 CP1 - critical journey and trace map (Register a loan) | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/1e96b9f |
| HU-06 | W11 CP2 - structured log contract | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/db2ed20 |
| HU-06 | W11 CP3 - traceparent middleware and JSON logger (circulation-api) | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/707ea43 |
| HU-06 | W11 CP3 - catalog client with traceparent and metrics registry | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/66a9f80 |
| HU-06 | W11 CP3 - context propagated to catalog-api, /metrics and test | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/374b96f |
| HU-06 | W11 CP3 - how to run with real console output | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/34a85fe |
| HU-06 | W11 CP4 - golden signals and metric catalog | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/58e1141 |
| HU-06 | W11 CP6 - SLIs, SLOs and error budget for loan registration | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/f92223a |
| HU-06 | W11 CP7 - burn-rate alerting rules and alert policy | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/c55faa4 |
| HU-06 | W11 CP8 - runbook for the fast-burn alert | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/1fd9248 |
| HU-06 | W11 CP9 - incident drill (slow catalog-api) and blameless postmortem | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/d9dd802 |
| HU-06 | W11 CP10 - hu-status for delivery 2 and implementation plan in the reflection document | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/11-week/workshop |

## 2. My individual contribution
- Chose HU-06 Register a loan as the critical journey: circulation-portal, gateway, circulation-api, membership-api, catalog-api and MongoDB loan_db.
- Wrote the hop table, Mermaid sequence diagram and span tree, marking where traceparent travels.
- Defined the JSON log contract with LMS domain fields and the list of data never logged (document number, email, JWT, internal token).
- Built a Go module for circulation-api with a traceparent middleware, a slog JSON logger, a catalog-api client that propagates the context, /metrics, and a test that checks the trace_id sent downstream.
- Wrote the golden-signals catalog for circulation-api, membership-api and catalog-api.
- Defined three SLOs for POST /api/v1/loans (availability 99.5 %, latency 99 % under 1 s, loan consistency 99.9 %) with the error budget and the team policy.
- Wrote multi-window burn-rate alerts validated with promtool, the alert policy and the runbook for the fast-burn alert.
- Ran an incident drill with a slow catalog-api and wrote a blameless postmortem with three user stories.

## 3. Blockers and risks
- The real membershipclient and catalogclient do not forward X-Correlation-Id or traceparent, so logs of one loan cannot be joined today.
- No Prometheus, Jaeger or log storage in our compose yet; the official gateway, lms-workflow and lms-worker are still empty.
- A timeout on loan-copy can leave a copy reserved without a loan (no saga compensation, ADR-008).

## 4. Plan for next week
- Upload lms-worker (due-date notices, fines accrual and catalog reindex jobs, ADR-009) to develop in small PRs, from the domain outwards, with unit tests.
- Propose to the team the postmortem stories: traceparent forwarding and a 1.5 s timeout with Idempotency-Key in the real catalogclient, and a separate connection pool for loan-copy in catalog-api.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Workshop folder: https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/11-week/workshop
- Reflection document: https://docs.google.com/document/d/1Xn1s14c0UCqSO5IJdMU5HpkndiJW7APpVQ5sLPenr0E/edit?usp=sharing
