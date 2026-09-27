<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       08-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 08

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Publish the REST API guidelines for 07-api and migrate the Membership bounded context (API and portal) out of the monolith into lms-membership-api and lms-membership-portal, following ADR-006, in small PRs that each build and pass checks on their own.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-API-004 | REST API guidelines for 07-api (versioning, pagination, error format, error-code catalogue, contract audit) | done | https://github.com/code-corhuila/library-docs/pull/19 |
| HU-02 | Student aggregate, Email value object and domain ports | done | https://github.com/code-corhuila/lms-membership-api/pull/3 |
| HU-02 | Create student use case | done | https://github.com/code-corhuila/lms-membership-api/pull/4 |
| HU-02, HU-03 | Infrastructure adapters: Postgres repository, circulation client, config and logger | done | https://github.com/code-corhuila/lms-membership-api/pull/5 |
| HU-02, HU-03 | HTTP response helpers and auth, correlation ID and CORS middleware | done | https://github.com/code-corhuila/lms-membership-api/pull/6 |
| HU-03 | Student management use cases (get, search, update, deactivate, suspend) | done | https://github.com/code-corhuila/lms-membership-api/pull/7 |
| HU-02, HU-03 | Health and student HTTP handlers and router | done | https://github.com/code-corhuila/lms-membership-api/pull/8 |
| HU-02, HU-03 | Service entry point, Dockerfile, Makefile and README | done | https://github.com/code-corhuila/lms-membership-api/pull/9 |
| HU-02, HU-03 | Portal scaffold (Vite, React, TypeScript, shared lib and UI) | done | https://github.com/code-corhuila/lms-membership-portal/pull/2 |
| HU-03 | Students list page | done | https://github.com/code-corhuila/lms-membership-portal/pull/3 |
| HU-02 | Student registration form page | done | https://github.com/code-corhuila/lms-membership-portal/pull/4 |
| HU-02, HU-03 | Portal routes and app entry point | done | https://github.com/code-corhuila/lms-membership-portal/pull/6 |

## 2. My individual contribution

**API guidelines — `code-corhuila/library-docs`**

- Wrote `07-api/guidelines.md`, a file the folder README listed as required and that did not exist. It closes eleven decisions (D-G01 to D-G11): versioning, pagination parameters and envelope, date ranges, error body and error-code catalogue, 403 in a single-role system, the 400/409/422 boundary, malformed UUIDs, missing 500 declarations and the two 401 codes.
- Audited the five existing OpenAPI contracts against those rules and recorded five defects (DF-1 to DF-5) and six gaps, each with an owner.
- Answered the review comments inside the PR and reconciled the overlap with a teammate's parallel PR on the same file. Merged as PR #19.

**Membership API migration — `code-corhuila/lms-membership-api`**

- Resumed the membership-service migration from the `lms-library` monolith. An earlier attempt had been opened as a single PR (#2) that did not build on its own and was closed without merging.
- Reordered the migration into seven slices following the hexagonal layers — domain, use cases, infrastructure adapters, HTTP middleware, handlers and router, entry point — so that every PR builds, passes `go vet` and runs its tests independently, and none exceeds the 400-line limit.
- Delivered each slice through its own `feat/HU-NN-...` or `chore/...` branch and a PR to `develop`, with Conventional Commits and implementation and tests committed separately.
- Kept the domain free of infrastructure: the `ActiveLoansChecker` port lets the deactivate use case ask circulation-service about active loans over HTTP instead of reading its database.
- Documented the known gaps in the PR descriptions instead of changing migrated code: the auth middleware does not pin the signing method and accepts tokens without `sub`, handlers have no tests yet, the repo has no `.gitignore`, `golang.org/x/crypto` is unused, and the container runs as root.

**Membership portal migration — `code-corhuila/lms-membership-portal`**

- Delivered the portal in four slices: project scaffold with shared lib and UI components, the students list page (HU-03), the registration form page (HU-02), and the routes and entry point that wire them together.
- Verified each slice with `tsc -b` and `oxlint`, and the last one with a full `npm run build`.
- Updated the README to correct the migration scope: the code comes from the `lms-library` frontend, it is not written from scratch.

## 3. Blockers and risks
- An open PR in the portal (#5) restructures it to follow the course's front-end repository norm (Anexo H) and adds the same entry files (`index.html`, `src/main.tsx`, `src/App.tsx`). Since #6 is already merged, #5 will need its conflicts resolved before it can merge.
- The membership API validates JWTs on every student route, including the two that circulation-service calls, so circulation-service will need a token signed with the same secret.
- The auth middleware issues found during the migration (no pinned signing method, tokens without `sub` accepted) are documented but not fixed yet.
- The portal duplicates `components/ui` and `lib` until the shared `lms-front` repository exists.

## 4. Plan for next week
- Coordinate with the author of #5 so the portal ends up aligned with Anexo H on top of the merged entry point.
- Fix the auth middleware: require HS256 and reject tokens without `sub`, with tests. Check whether `lms-access-api` has the same issue.
- Add tests for the membership HTTP handlers and router.
- Repo hygiene PR for the API: `.gitignore`, `go mod tidy`, non-root container user and a `.gitattributes` to normalize line endings.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Membership API pull requests: https://github.com/code-corhuila/lms-membership-api/pulls?q=is%3Apr+author%3AHermesss29
- Membership portal pull requests: https://github.com/code-corhuila/lms-membership-portal/pulls?q=is%3Apr+author%3AHermesss29
- API guidelines (merged): https://github.com/code-corhuila/library-docs/pull/19
- Portal entry point: https://github.com/code-corhuila/lms-membership-portal/pull/6
