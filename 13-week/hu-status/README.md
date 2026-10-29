<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Threat model the trust boundary circulation-api -> membership-api in the W13 workshop (attack surface, STRIDE PR, security backlog and security gates in CI) and start uploading lms-worker for HU-08.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-08 | W13 CP0 - claim trust boundary circulation-api -> membership-api | done | https://github.com/code-corhuila/library-docs/issues/33 |
| HU-08 | W13 CP1 - attack surface and trust boundaries | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/a2e75f1 |
| HU-08 | W13 CP2 - STRIDE on circulation-api -> membership-api (PR to library-docs) | doing | https://github.com/code-corhuila/library-docs/pull/40 |
| HU-03 | W13 CP3 - [security] restrict service tokens to the routes circulation-api needs (E, P1) | todo | https://github.com/code-corhuila/library-docs/issues/41 |
| HU-06 | W13 CP3 - [security] verify sub, aud and exp of service tokens (S, P1) | todo | https://github.com/code-corhuila/library-docs/issues/42 |
| HU-06 | W13 CP3 - [security] return only eligibility data to circulation-api (I, P1) | todo | https://github.com/code-corhuila/library-docs/issues/43 |
| HU-03 | W13 CP4 - security gates and negative authorization test | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/commit/dbfbff0 |
| HU-08 | lms-worker scaffold: module, config and logger | doing | https://github.com/code-corhuila/lms-worker/pull/2 |

## 2. My individual contribution
- Claimed and modeled the trust boundary circulation-api -> membership-api (GET /api/v1/students/{id}, POST /api/v1/students/{id}/suspend, HS256 service token).
- Drew the data flow diagram of the whole LMS with trust boundaries TB1 to TB5, an asset inventory and the entry points of my boundary.
- Wrote the STRIDE table with six concrete threats from the real code (RequireAuth accepts any service token for the whole /students group, unvalidated suspension days, full PII returned to circulation), risk ranking and residual risk, citing Pressman (7th ed., pp. 318, 378, 454) from the CORHUILA library.
- Turned the three highest risks into security stories #41, #42 and #43 with verifiable acceptance criteria, labels security and P1.
- Proposed a security workflow (gitleaks, govulncheck, negative authorization test) and a Go test that proves a service token cannot deactivate a student.
- Opened lms-worker PR #2 with the module, configuration and logger.
- Peer review: https://github.com/code-corhuila/library-docs/pull/36#pullrequestreview-5699486741 (row T of the circulation-api -> catalog-api STRIDE, PR #36)

## 3. Blockers and risks
- The negative authorization test fails today on purpose: route-level authorization for service tokens does not exist yet (story #41).
- One INTERNAL_JWT_SECRET is shared by every Go service; changing the key model needs a team decision.
- The security workflow is a proposal in my fork; it still has to be added to lms-membership-api.

## 4. Plan for next week
- Answer the review on PR #40, finish the STRIDE and refine stories #41-#43 (delivery 2).
- Run the negative test against lms-membership-api and record the real output.
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
