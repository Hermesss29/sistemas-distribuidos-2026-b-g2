<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       09-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 09

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Deliver the Circulation portal (HU-06, HU-07, HU-08) as a Module Federation remote of lms-front in small PRs that each build on their own, and promote the Membership context (db, api, portal) from develop to qa following the re-application policy.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-06, HU-07, HU-08 | Circulation portal scaffold as a Module Federation remote, CI and placeholder routes | done | https://github.com/code-corhuila/lms-circulation-portal/pull/2 |
| HU-07, HU-08 | Loans list with return action and overdue loans report | done | https://github.com/code-corhuila/lms-circulation-portal/pull/3 |
| HU-06 | Loan registration form with idempotency key, searchable student/book select, deploy files | done | https://github.com/code-corhuila/lms-circulation-portal/pull/4 |
| HU-02, HU-03 | Fix gofmt in the student handler so the API pipeline can pass | done | https://github.com/code-corhuila/lms-membership-api/pull/12 |
| HU-02, HU-03 | Promote membership database to qa | doing (in review) | https://github.com/code-corhuila/lms-membership-db/pull/6 |
| HU-02, HU-03 | Promote membership API to qa | doing (in review) | https://github.com/code-corhuila/lms-membership-api/pull/13 |
| HU-02, HU-03 | Promote membership portal to qa | doing (in review) | https://github.com/code-corhuila/lms-membership-portal/pull/7 |

## 2. My individual contribution

**Circulation portal — `code-corhuila/lms-circulation-portal`**

- Delivered the portal in three PRs, each building and linting on its own (`npm run build`, `oxlint`) and none over the 400-line limit:
  - **#2:** the scaffold. The portal plugs into `lms-front` as a Module Federation remote (Annex H), with no HTTP client or session of its own; it uses `shell/apiClient` and `shell/session`. Includes CI and the PR template. The three screens start as placeholders.
  - **#3:** the loans list with a status filter and return action (HU-07) and the read-only overdue report (HU-08). Both have the four screen states: loading, error with retry, empty and data.
  - **#4:** the loan registration form (HU-06), which sends an `Idempotency-Key` so a retry after a network cut cannot reserve two copies, plus the searchable student/book select and the Docker/nginx/compose files.
- Wrote the placeholder routes so the portal compiled from the first PR, and replaced them one HU at a time.
- Fixed real issues found while delivering:
  - A stale contract reference in the types (`library-api.yaml` no longer exists).
  - An nginx rule that never applied because it pointed `remoteEntry.js` to the wrong path.
  - Outdated "known gaps" in the README.
  - Kept the plugin's `.mf/` cache (with a local machine path inside) out of the repo.
- Documented what is still missing instead of hiding it: the searchable select shows "No matches" when the search fails, with no error or retry state.

**Membership API — `code-corhuila/lms-membership-api`**

- **#12:** found that `develop` did not pass the `gofmt` step of its own `ci.yml` (unsorted imports introduced by an earlier PR). Fixed it before promoting, so `qa` would not receive a failing pipeline.

**Promotion of the Membership context to `qa` (db, api, portal)**

- Promoted `develop` to `qa` in the three repos following `branching-policy.md`: no merge between permanent branches, a branch cut from `qa`, and each PR already merged in `develop` re-applied with `git cherry-pick -x` (`-m 1` for merge commits) in first-parent order: 4 commits in db, 11 in api and 6 in portal.
- One commit per PR instead of every individual commit. In the portal, one merge in `develop` held a conflict resolution between two PRs that create the same files; re-applying the loose commits would have lost it.
- Verified each promotion before opening it:
  - the resulting tree is identical to `develop` (empty diff);
  - the commit count matches the number of `(cherry picked from commit …)` lines;
  - the build passes (`go build/vet/test` + `gofmt` for the API, `npm ci` + `npm run build` for the portal).
  Every PR lists the full promotion trail.
- Found that the policy's `qa/...` branch name cannot exist in git while a `qa` branch exists (`refs/heads/qa` blocks `refs/heads/qa/*`). Used `promote-qa/...`, the prefix the team already uses in the circulation repos, and explained it in each PR.

**Retrospective forum (Corte 2)**

- Wrote the team retrospective on the membership migration: the root cause (slices cut by file type instead of by dependency order), the solution and its verification.
- Replied to another team's post on a Maven dependency-mediation bug with three concrete additions:
  - the same trap can come back through the parent POM's `dependencyManagement`;
  - `dependency:tree -Dverbose` shows the conflict directly;
  - a compose `healthcheck` with `up --wait` makes the pipeline test the packaged artifact, not just the classes.

## 3. Blockers and risks
- **CI blocked:** GitHub Actions is blocked by the `code-corhuila` organization's billing. Every check this week was run locally with the same steps the workflows run.
- **No contract tests:** `qa` requires contract tests (Pact) as its gate, and no repo has them yet. The promotions document this instead of skipping it silently.
- **Policy example unusable:** the `qa/...` branch example in `branching-policy.md` cannot be created in git. Until the policy is updated, the team uses `promote-qa/...`.
- **Circulation portal not tested end to end:** it builds and lints, but it has not run against a real `lms-front` host yet.
- **Known findings in develop:** `lms-membership-portal` reports 4 npm audit vulnerabilities (2 moderate, 2 high).

## 4. Plan for next week
- Get the three `qa` promotion PRs (db #6, api #13, portal #7) reviewed and merged.
- Fix the searchable select's error state (show the error with a retry instead of "No matches"), closing the documented gap in the circulation portal.
- Open a PR to `develop` in `lms-membership-portal` for the npm audit findings.
- Add tests for the membership HTTP handlers and router.
- Propose in `library-docs` that the promotion branch prefix be `promote-qa/` instead of `qa/`.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Circulation portal scaffold: https://github.com/code-corhuila/lms-circulation-portal/pull/2
- Loans list and overdue report (HU-07, HU-08): https://github.com/code-corhuila/lms-circulation-portal/pull/3
- Loan registration form (HU-06): https://github.com/code-corhuila/lms-circulation-portal/pull/4
- gofmt fix before promotion: https://github.com/code-corhuila/lms-membership-api/pull/12
- Promotion to qa — db: https://github.com/code-corhuila/lms-membership-db/pull/6
- Promotion to qa — api: https://github.com/code-corhuila/lms-membership-api/pull/13
- Promotion to qa — portal: https://github.com/code-corhuila/lms-membership-portal/pull/7
