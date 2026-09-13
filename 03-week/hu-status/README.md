<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       03-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 03

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Define the product backlog (epics and user stories) and the non-functional requirements for LMS-LIBRARY-V1.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01..HU-08 | Backlog definition (4 epics, 8 user stories) | done | [03-week/hu-status/HU-BACKLOG.md](./HU-BACKLOG.md) |

## 2. My individual contribution
- Authored [HU-BACKLOG.md](./HU-BACKLOG.md): 4 epics and 8 user stories (HU-01 to HU-08), each with verifiable acceptance criteria, MoSCoW priority, and dependencies; plus the Definition of Ready and Definition of Done. Work session held on 2026-08-22.
- Authored non-functional requirements in `04-requirements/non-functional.md` with justified metrics (P95 < 300ms, P99 < 600ms under 50 RPS, minimum throughput of 20 RPS, 95% availability SLO with justification, limit of 3 instances per service in V1) and critical endpoints linked to FR-004, FR-010 and FR-015.
  Evidence: https://github.com/code-corhuila/library-docs/commit/0aa379e
- Participated in decisions for the week: fixed 7-day loan policy with no renewals, maximum of 2 active loans per student, 7-day suspension for late returns, and a single Administrator operating role.

## 3. Blockers and risks
- Non-functional targets (latency, throughput, availability) still need validation against the chosen infrastructure once the MVP architecture is defined.

## 4. Plan for next week
- Start development of the first prioritized user stories from the backlog.
- Migrate and refine the backlog into the official documentation repository.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [ ] No secrets; config via environment variables

## 6. Evidence links
- Backlog: [03-week/hu-status/HU-BACKLOG.md](./HU-BACKLOG.md) (commit 2026-08-23, this repo)
- Non-functional requirements: https://github.com/code-corhuila/library-docs/commit/0aa379e
