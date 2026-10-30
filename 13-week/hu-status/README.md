<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Threat model the trust boundary circulation-api -> membership-api in the W13 workshop (attack surface, STRIDE PR answered after peer review, refined security backlog and a negative authorization test run against the real membership-api) and start uploading lms-worker for HU-08.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-08 | W13 CP0 - claim trust boundary circulation-api -> membership-api | done | https://github.com/code-corhuila/library-docs/issues/33 |
| HU-08 | W13 CP1 - attack surface and trust boundaries | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/a2e75f1 |
| HU-08 | W13 CP2 - STRIDE on circulation-api -> membership-api (PR to library-docs) | done | https://github.com/code-corhuila/library-docs/pull/40 |
| HU-03 | W13 CP3 - [security] restrict service tokens to the routes circulation-api needs (E, P1) | todo | https://github.com/code-corhuila/library-docs/issues/41 |
| HU-06 | W13 CP3 - [security] verify sub, aud and exp of service tokens (S, P1) | todo | https://github.com/code-corhuila/library-docs/issues/42 |
| HU-06 | W13 CP3 - [security] return only eligibility data to circulation-api (I, P1) | todo | https://github.com/code-corhuila/library-docs/issues/43 |
| HU-03 | W13 CP4 - security gates and negative authorization test | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/dbfbff0 |
| HU-08 | W13 CP6 - review on PR #40 answered: replayed suspend (T) and loan_id audit (R) | done | https://github.com/code-corhuila/library-docs/commit/d37b8a1 |
| HU-03 | W13 CP7 - security stories refined and linked to #37 and #39 | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/7bdd3aa |
| HU-03 | W13 CP8 - negative authorization test run on membership-api (FAIL as expected, 200 instead of 403) | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/fdb75ff |
| HU-08 | lms-worker scaffold: module, config and logger | doing | https://github.com/code-corhuila/lms-worker/pull/2 |

## 2. My individual contribution
- Claimed and modeled the trust boundary circulation-api -> membership-api (GET /api/v1/students/{id}, POST /api/v1/students/{id}/suspend, HS256 service token).
- Drew the data flow diagram of the whole LMS with trust boundaries TB1 to TB5, an asset inventory and the entry points of my boundary.
- Wrote the STRIDE table with six concrete threats from the real code (RequireAuth accepts any service token for the whole /students group, unvalidated suspension days, full PII returned to circulation), risk ranking and residual risk, citing Pressman (7th ed., pp. 318, 378, 454) from the CORHUILA library.
- Turned the three highest risks into security stories #41, #42 and #43 with verifiable acceptance criteria, labels security and P1.
- Proposed a security workflow (gitleaks, govulncheck, negative authorization test) and a Go test that proves a service token cannot deactivate a student.
- Answered the peer review on PR #40: added the replayed suspend to row T (Idempotency-Key per loan and GREATEST on suspended_until) and changed the audit in row R to log loan_id instead of the shared-secret sub.
- Linked my stories to the neighbour boundary: #41 with #37 and #42 with #39 (same RequireAuth gap in catalog-api).
- Ran the negative authorization test against RequireAuth of lms-membership-api (develop): the service token reaches deactivate, edit and list with 200, which proves story #41. security.yml validated with actionlint 1.7.7.
- Opened lms-worker PR #2 with the module, configuration and logger.
- Peer review: https://github.com/code-corhuila/library-docs/pull/36#pullrequestreview-5699486741 (row T of the circulation-api -> catalog-api STRIDE, PR #36)

## 3. Blockers and risks
- The negative authorization test fails today on purpose: route-level authorization for service tokens does not exist yet (story #41).
- One INTERNAL_JWT_SECRET is shared by every Go service; changing the key model needs a team decision.
- The security workflow is a proposal in my fork; it still has to be added to lms-membership-api.

## 4. Plan for next week
- Propose story #41 (route-level authorization for service tokens) for the next sprint, so the negative test turns green.
- Agree with the catalog owner on one claim set (sub, aud, exp, scope) for #42 and #39.
- Continue lms-worker: domain ports, use cases with tests, HTTP clients and scheduler.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Workshop folder: https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/13-week/workshop
- STRIDE PR: https://github.com/code-corhuila/library-docs/pull/40
- Reflection document: https://docs.google.com/document/d/1lgmE3VAZHeHtZZbuf14te5gzRz_AdVDH8Z5tWrZ_7tU/edit?usp=sharing
