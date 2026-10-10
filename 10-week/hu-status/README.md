<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       10-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 10

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Release v2.0.1 of the Membership context (db, api, portal) and of lms-front from qa to main following the course release flow (release branch from main, HU branch re-applied with cherry-pick -x, teacher-approved PR, annotated tag), and answer the automated review findings.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-02, HU-03 | Promote membership database, API and portal to qa (merged this week) | done | https://github.com/code-corhuila/lms-membership-db/pull/6 |
| HU-02, HU-03 | Membership database: HU branch to release/2.0.1 | done | https://github.com/code-corhuila/lms-membership-db/pull/7 |
| HU-02, HU-03 | Membership database: release v2.0.1 to main + tag | done | https://github.com/code-corhuila/lms-membership-db/pull/8 |
| HU-02, HU-03 | Membership API: HU branch to release/2.0.1 | done | https://github.com/code-corhuila/lms-membership-api/pull/14 |
| HU-02, HU-03 | Membership API: release v2.0.1 to main + tag | done | https://github.com/code-corhuila/lms-membership-api/pull/15 |
| HU-02, HU-03 | Membership portal: HU branch to release/2.0.1 | done | https://github.com/code-corhuila/lms-membership-portal/pull/8 |
| HU-02, HU-03 | Membership portal: release v2.0.1 to main + tag | done | https://github.com/code-corhuila/lms-membership-portal/pull/9 |
| HU-02 | Track the non-transactional idempotency save found in review | todo (issue open) | https://github.com/code-corhuila/lms-membership-api/issues/16 |
| HU-01 | lms-front: bring the stress/perf tests to develop so develop and qa match | done | https://github.com/code-corhuila/lms-front/pull/13 |
| HU-01 | lms-front: HU branch to release/2.0.1 | done | https://github.com/code-corhuila/lms-front/pull/14 |
| HU-01 | lms-front: release v2.0.1 to main + tag | done | https://github.com/code-corhuila/lms-front/pull/15 |

## 2. My individual contribution

**Release v2.0.1 of the Membership context — db, api, portal**

- Got the three `qa` promotions from last week (db #6, api #13, portal #7) reviewed and merged.
- Ran the course release flow in each repo:
  - cut `release/2.0.1` from `main`;
  - opened one HU branch (`hu-02-03-release`) that re-applies the `qa` commits with `git cherry-pick -x` (4 in db, 11 in api, 6 in portal), so every commit carries two trailers: one from `develop` and one from `qa`;
  - opened the PR `release/2.0.1` → `main`, approved by the teacher (`ariel5253`) and merged with a merge commit;
  - created the annotated tag `v2.0.1` on that merge commit (db `33a74cd`, api `4181abd`, portal `f62ed7f`).
- Checked every hop before opening it:
  - the tree is identical to `qa` (empty `git diff`);
  - the commit count matches the number of trailers;
  - the build passes;
  - the merge into `main` has two parents (merge commit, not squash).
- Caught that `git switch -c X origin/Y` leaves the new branch tracking `origin/Y`, so a plain push could land in the wrong branch. I always pushed with `git push -u origin <branch>` and checked `@{u}` afterwards.

**Answering the automated review (db #8)**

- Wrote a decision for each finding (applied and how, or not applied and why): https://github.com/code-corhuila/lms-membership-db/pull/8#issuecomment-6040735567
  - Changed the title to `chore(release): …`.
  - Replaced the old "Actions blocked" note with the real `db-ci` run (green) in the description.
  - Proved the merge commit in the trail is by design: `git diff` between the merge and its second parent is empty.
  - Explained that the 400-line cap applies to PRs to `develop`.
- For the foreign-key / idempotency finding, I confirmed in the code that the cause is in the API, not in the schema: `Students.Save` and `Idempotency.Save` are two separate `Exec` calls with no transaction. Opened api #16 with the failure scenario and a proposed fix, and left the fix for a later release instead of changing code inside a release.

**Tag push failure and workaround**

- `git push origin refs/tags/v2.0.1` was rejected twice with a GitHub `Internal Server Error`. Before retrying anything else, I checked that:
  - nothing had been created;
  - no branch had moved;
  - no ruleset targets tags.
- Created the same annotated tag through the Git Data API (`git/tags` + `git/refs`), on the same commit and with the same message, and verified it with `git ls-remote`. Used the same method for api, portal and lms-front, so each repo has exactly one tag.

**lms-front release v2.0.1**

- Audited the clone before releasing and found that `qa` held a commit that `develop` did not (`a5f7c73`, connection/stress/performance tests).
  - Cause: #12 was merged into `test/shell-unit-tests` instead of `develop`. #9 was merged without deleting its branch, so GitHub never retargeted #12.
  - PR #10 itself said that `qa` must never hold something `develop` does not.
- Opened **#13**, which brings the same commit (`a76845d`) to `develop`. Merged it with a merge commit so the hash quoted in the `qa` trailer exists in `develop`. After it, `git diff develop qa` is empty.
- Released through `release/2.0.1`:
  - **#14:** 8 commits, 16 trailers.
  - **#15** → `main`: approved by `ariel5253` and `OscarAreiza`.
  - Tag `v2.0.1` on `b012ea2`.
- The repo has no CI, so I ran the checks locally and listed the results in the PRs:
  - `npm ci` and `npm run build` pass;
  - `oxlint` is clean;
  - `npm test` 65/65;
  - `npm run test:stress` 8/8.

## 3. Blockers and risks
- **Critical finding in lms-front `main` (#15 review):** `src/bootstrap.tsx` accepts a `devToken` query parameter and is not gated by `import.meta.env.DEV`. Any URL with `?devToken=x` gets past `RequireAuth`. It came with the shell migration and is now in `v2.0.1`; it has to be fixed first, through `develop` → `qa` → a new release.
- **Review findings still without a written answer:** api #15 (routes not scoped by token type; idempotency race, tracked in #16), portal #9 (shell remote URL hardcoded to `localhost:3000`, no tests, CI without lint), lms-front #15.
- **Known vulnerabilities promoted to `main` with no tracking issue:**
  - membership portal: 2 moderate, 2 high;
  - lms-front: 2 moderate, 3 high.
- **Tag push through git fails in this org** with `Internal Server Error`, with no rule behind it. The API works, but it is a workaround, not a fix.
- **No contract tests (Pact) in any repo, and lms-front has no CI.** The checks are local only.
- **Standards gaps flagged by the review:** the api #15 and portal #9 titles use `release:`, which is not a Conventional Commits type (db #8 was corrected). The standard does not say whether the 400-line cap applies to release PRs.

## 4. Plan for next week
- Answer each finding in api #15, portal #9 and lms-front #15 the same way as in db #8, opening an issue for everything that is deferred.
- Remove or gate the `devToken` bypass in lms-front (PR to `develop`, then `qa` and a new release).
- Fix api #16: save the student and the idempotency key in one transaction, with a test for the retry race.
- Open issues and PRs to `develop` for the npm audit findings in membership portal and lms-front.
- Make the shell remote URL of the membership portal configurable per environment.
- Propose in `library-docs`:
  - the `promote-qa/` prefix;
  - the line-cap rule for release PRs;
  - the API method for tags.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Membership db — qa promotion: https://github.com/code-corhuila/lms-membership-db/pull/6
- Membership db — release PR to main: https://github.com/code-corhuila/lms-membership-db/pull/8
- Membership db — answer to the review findings: https://github.com/code-corhuila/lms-membership-db/pull/8#issuecomment-6040735567
- Membership db — tag: https://github.com/code-corhuila/lms-membership-db/tree/v2.0.1
- Membership API — qa promotion: https://github.com/code-corhuila/lms-membership-api/pull/13
- Membership API — release PR to main: https://github.com/code-corhuila/lms-membership-api/pull/15
- Membership API — tag: https://github.com/code-corhuila/lms-membership-api/tree/v2.0.1
- Membership API — idempotency transaction issue: https://github.com/code-corhuila/lms-membership-api/issues/16
- Membership portal — qa promotion: https://github.com/code-corhuila/lms-membership-portal/pull/7
- Membership portal — release PR to main: https://github.com/code-corhuila/lms-membership-portal/pull/9
- Membership portal — tag: https://github.com/code-corhuila/lms-membership-portal/tree/v2.0.1
- lms-front — develop aligned with qa: https://github.com/code-corhuila/lms-front/pull/13
- lms-front — HU branch to release: https://github.com/code-corhuila/lms-front/pull/14
- lms-front — release PR to main: https://github.com/code-corhuila/lms-front/pull/15
- lms-front — tag: https://github.com/code-corhuila/lms-front/tree/v2.0.1
