<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       15-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 15

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Close MVP3: cut the v3.0.0 release of lms-worker (HU-08) and lms-access-db (HU-01) from qa to main, and present the first part of the final project defense (c3-release, moment 1) with the whole LMS running locally and a controlled failure.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-08 | lms-worker release v3.0.0 to main (cherry-pick -x from qa, 24 commits with develop and qa trail) | doing | https://github.com/code-corhuila/lms-worker/pull/8 |
| HU-01 | lms-access-db release v3.0.0 to main (cherry-pick -x from qa, 16 commits with develop and qa trail) | doing | https://github.com/code-corhuila/lms-access-db/pull/8 |
| HU-01, HU-08 | Final project defense, part 1 (c3-release moment 1): end-to-end demo, controlled failure and individual questions | done | In-class defense, c3-release rubric (S1, S2, S3) |
| HU-01..HU-09 | Local full-stack environment for the demo: 17 repos from develop, shared PostgreSQL and MongoDB (Anexo J), 6 Liquibase migrations, 5 APIs, workflow, worker and gateway | done | Local environment (lms-qa-2), built from develop on 2026-11-11 |

## 2. My individual contribution
- Opened the v3.0.0 release PRs to main for lms-worker and lms-access-db. They were cut from main and filled with git cherry-pick -x from qa, so every commit carries both the develop and the qa "cherry picked from" trail. ariel5253 is requested as reviewer.
- Verified the release content before opening the PRs: git diff against qa is empty in both repos, and go build, go vet and go test pass on the worker release branch.
- Presented the first part of the final defense (c3-release, moment 1): the LMS ran end to end locally, with a controlled failure (stopping a service while the rest kept working), and answered questions about my own work and about the project's technical decisions.
- Built a new local environment (lms-qa-2) that runs the whole LMS from develop: 17 repos, one PostgreSQL and one MongoDB instance with one schema and one login user per domain (Anexo J), migrations for access, membership, catalog, iam, workflow and loan_db run in order, access, membership, catalog, circulation and iam APIs, lms-workflow, lms-worker and a gateway based on lms-api-gateway (ADR-007 routing and rate limits plus CORS for the dev front). It includes a QA administrator seed and a check script for health, login and one read per domain.
- Reviewed library-docs against the code and documented for the team the inconsistencies found: loan period 3 days in code and in ADR-009 vs 7 days in older documents, the ADR register stopping at ADR-010 while records go up to ADR-013, and repo-migration-map still marking lms-worker and lms-access-db as not started.

## 3. Blockers and risks
- Both v3.0.0 release PRs are waiting for the teacher's approval on main, so the v3.0.0 tag cannot be created yet.
- lms-worker has unit tests but no CI workflow yet, and neither lms-worker nor lms-access-db has a CHANGELOG.md.
- The worker and the workflow call routes that lms-circulation-api and lms-catalog-api do not expose yet (/loans/{id}/notify, /loans/accrue-fines, /books/reindex, /loans/{id}/mark-late), so those jobs log 404 failures (also recorded in ADR-013).

## 4. Plan for next week
- Merge the v3.0.0 release PRs as soon as they are approved and create the annotated v3.0.0 tag on main in both repos, before the code deadline (2026-11-15 23:59).
- Add the CI workflow (go vet, go test) to lms-worker and a CHANGELOG.md to lms-worker and lms-access-db, promoted to main through the branching policy.
- Update library-docs: repo-migration-map statuses and the ADR register (ADR-011 to ADR-013).

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- lms-worker release v3.0.0: https://github.com/code-corhuila/lms-worker/pull/8
- lms-access-db release v3.0.0: https://github.com/code-corhuila/lms-access-db/pull/8
- Week 14 status (develop and qa delivery): https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/blob/main/14-week/hu-status/README.md
