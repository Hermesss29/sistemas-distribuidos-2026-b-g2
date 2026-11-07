<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       14-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 14

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Finish uploading lms-worker (HU-08) and lms-access-db (HU-01) to develop with tests and documentation, verify the access-db migrations locally with Docker, and promote both repositories to qa by re-application (git cherry-pick -x).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-08 | lms-worker scaffold: module, config and logger | done | https://github.com/code-corhuila/lms-worker/pull/2 |
| HU-08 | lms-worker driven ports and the three jobs (due-date notices, fines accrual, catalog reindex) with unit tests | done | https://github.com/code-corhuila/lms-worker/pull/3 |
| HU-08 | lms-worker HTTP clients, internal HS256 service token and job ticker with tests | done | https://github.com/code-corhuila/lms-worker/pull/4 |
| HU-08 | lms-worker entry point and multi-stage Dockerfile | done | https://github.com/code-corhuila/lms-worker/pull/5 |
| HU-08 | lms-worker README: structure, configuration, run/test/build and known integration gap | done | https://github.com/code-corhuila/lms-worker/pull/6 |
| HU-01 | lms-access-db schema access, table administrator and username index with rollbacks | done | https://github.com/code-corhuila/lms-access-db/pull/2 |
| HU-01 | lms-access-db roles and grants, changelog master, deploy executor and reconstruction CI | done | https://github.com/code-corhuila/lms-access-db/pull/3 |
| HU-01 | lms-access-db local reconstruction check with Docker (update, no-op update, rollback-count 999, update) | done | https://github.com/code-corhuila/lms-access-db/pull/3#issuecomment-6677534241 |
| HU-01 | lms-access-db README: structure, access_app user and local verification | done | https://github.com/code-corhuila/lms-access-db/pull/4 |
| HU-01 | lms-access-db db-ci loads deploy/liquibase.properties and checks the per-domain changelog tables; .env.example aligned with compose.yml | done | https://github.com/code-corhuila/lms-access-db/pull/6 |
| HU-08 | Promote lms-worker to qa (cherry-pick -x) | done | https://github.com/code-corhuila/lms-worker/pull/7 |
| HU-01 | Promote lms-access-db to qa (cherry-pick -x) | done | https://github.com/code-corhuila/lms-access-db/pull/5 |
| HU-01 | Promote the lms-access-db CI fix to qa (cherry-pick -x) | done | https://github.com/code-corhuila/lms-access-db/pull/7 |

## 2. My individual contribution
- Finished lms-worker in small PRs to develop. It is a stateless Go service with hexagonal layers: driven ports in domain/notice, the three scheduled jobs in application/usecase and the HTTP clients, ticker and logger in infrastructure.
- Kept the worker polling-based as ADR-009 decides: an internal ticker runs every job every POLL_INTERVAL (15m by default). There is no broker (AT-002), and the README now corrects the earlier "event-driven" wording.
- Added unit tests for every job, the service token, both HTTP clients and the ticker. A failing job is logged and does not stop the others.
- Documented the integration gap: the worker calls POST /loans/{id}/notify, POST /loans/accrue-fines and POST /books/reindex, which lms-circulation-api and lms-catalog-api do not expose yet.
- Uploaded lms-access-db with Liquibase 4.29 following Anexo A: schema access, table administrator (uuid, text with CHECK, unique username index), NOLOGIN roles access_reader/access_writer, the grant to access_app (Anexo J, J.7) and one rollback per changeset.
- Found that the migration grants a role to access_app, a login that only lms-infra-postgres creates, so it fails on an empty CI database. Added a CI step that creates access_app first (57c64fa).
- The org CI was blocked by a billing issue, so I ran the same reconstruction check locally with Docker. Without access_app it fails at dcl-grants-002 (Run: 6 of 7). With the fix the four steps pass. I left the evidence as a comment on PR #3.
- Wrote both READMEs and promoted both repositories from develop to qa by re-application (git cherry-pick -x), never by merge.
- Closed the follow-up from the retrospective: db-ci and the local check now load deploy/liquibase.properties, so migrations are tracked in databasechangelog_access as on the shared instance, and a CI step checks those tables exist. .env.example now lists the variables compose.yml reads.
- Wrote the retrospective forum post about the access_app failure and replied to a teammate's post on the idempotency_keys upsert grant in lms-circulation-db.

## 3. Blockers and risks
- The org account is locked by a billing issue, so db-ci has not run on GitHub yet. The local Docker check covers it until billing is fixed and the job is re-run.
- The worker endpoints listed in its README do not exist yet in lms-circulation-api and lms-catalog-api. Until they do, those jobs log a failure on every tick.

## 4. Plan for next week
- Cut the MVP3 release from main (release/<version>, cherry-pick -x from qa) once the version is announced, get the approval from ariel5253 and tag main.
- Agree with the circulation and catalog owners on the three endpoints the worker needs.
- Re-run db-ci on GitHub when billing is fixed.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- lms-worker develop: https://github.com/code-corhuila/lms-worker/tree/develop
- lms-access-db develop: https://github.com/code-corhuila/lms-access-db/tree/develop
- Local reconstruction check (Docker): https://github.com/code-corhuila/lms-access-db/pull/3#issuecomment-6677534241
- Promotion to qa: https://github.com/code-corhuila/lms-worker/pull/7 and https://github.com/code-corhuila/lms-access-db/pull/5
- Local check with liquibase.properties: https://github.com/code-corhuila/lms-access-db/pull/6
