<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       04-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 04

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Develop HU-03 (student management) and support peer review of teammates' user stories.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-03 | Student search, editing, and deactivation | doing (in review) | https://github.com/OscarAreiza/lms-library/pull/3 |

## 2. My individual contribution
- Migrated and refined the backlog into the official documentation repository: https://github.com/code-corhuila/library-docs/commit/cae4dd2 (2026-08-24, `04-requirements/user-stories.md`).
- Created Jira card LMS-16 "HU-03-StudentManagement", assigned to myself, on 2026-08-26; moved to "In Review" on 2026-08-28.
  https://lms-library.atlassian.net/browse/LMS-16
- Opened Pull Request #3 in OscarAreiza/lms-library on 2026-08-27: "feat: implement HU-03 student search, editing, and deactivation".
  https://github.com/OscarAreiza/lms-library/pull/3
- Reviewed the PR "feat(catalog): register a new book (HU-04)" on 2026-08-27.
- Reviewed the PR "Feat/hu 04 book registration" on 2026-08-30.
- Scope of HU-03: student search and listing, editing of contact data, and blocking deactivation when a student has active loans or an active suspension.

## 3. Blockers and risks
- PR #3 was closed without being merged; HU-03 required further work before it could be integrated.

## 4. Plan for next week
- Address review feedback on HU-03 and prepare it for a new merge attempt.
- Continue peer reviews for teammates' pull requests.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [ ] No secrets; config via environment variables

## 6. Evidence links
- https://github.com/code-corhuila/library-docs/commit/cae4dd2
- https://lms-library.atlassian.net/browse/LMS-16
- https://github.com/OscarAreiza/lms-library/pull/3
