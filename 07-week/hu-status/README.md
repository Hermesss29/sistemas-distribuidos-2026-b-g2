<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       07-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 07

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Migrate the Access bounded context out of the monolith into its own repository (lms-access-api) following ADR-006 repo-per-context decomposition, and document the architectural decisions that the new structure requires.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01 | Login use case migrated to access-service (credential validation and token issuing) | done | https://github.com/code-corhuila/lms-access-api/pull/3 |
| HU-01 | Administrator aggregate migrated to access-service (identity, credential validation, invariants) | done | https://github.com/code-corhuila/lms-access-api/pull/4 |
| HU-ARCH-007 | ADR-007 — gateway auth boundary and rate limiting | doing (in review) | https://github.com/code-corhuila/library-docs/pull/8 |
| HU-ARCH-009 | ADR-009 — worker scheduling model (polling vs. message broker) | doing (in review) | https://github.com/code-corhuila/library-docs/pull/10 |

## 2. My individual contribution

**Access service migration — `code-corhuila/lms-access-api`**

- Ported the Login use case (`internal/config/application/usecase/login.go`) from the monolith, implementing the HU-01 login flow: credential validation and token issuing. Delivered through a feature branch and Pull Request #3, merged after review.
- Ported the Administrator aggregate (`internal/config/domain/access/administrator.go`): identity, credential validation and the invariants of the Access bounded context. Delivered through Pull Request #4, merged after review.
- Added the access domain ports, keeping the domain free of infrastructure dependencies.
- Wrote unit tests for both slices — Login use case and Administrator aggregate — committed separately from the implementation.

**Architecture documentation — `code-corhuila/library-docs`**

- Authored ADR-007, defining the auth boundary of `lms-api-gateway`: authentication stays at routing level in the gateway while JWT validation remains inside each backend service, consistent with ADR-003.
- Authored ADR-009, documenting that `lms-worker` uses a polling scheduling model instead of a message broker or event bus, consistent with the architectural debt already accepted in AT-002.

**Peer review**

- Reviewed `lms-access-api` PR #2 (entry point and tooling migration) and PR #5 (infrastructure adapters).
- Reviewed `library-docs` PR #5 (security rules), PR #7 (ADR-006 repo-per-context decomposition), PR #9 (ADR-008 workflow saga scoping) and PR #12 (git conventions alignment).

## 3. Blockers and risks
- ADR-007 and ADR-009 are still open and pending team review; the decisions they document are being applied in code before being formally accepted.
- The infrastructure adapters for the access service (PR #5) are still in review, so the service is not yet runnable end to end in its new repository.
- Only the Access context has been migrated. Membership, Catalog and Circulation still live in the monolith, so the system currently runs in two places at once.

## 4. Plan for next week
- Finish the access-service migration: merge the infrastructure adapters and wire the service behind the API gateway.
- Get ADR-007 and ADR-009 reviewed and merged.
- Start the migration of the next bounded context following the same slice-per-PR approach.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Login use case migration (merged): https://github.com/code-corhuila/lms-access-api/pull/3
- Administrator aggregate migration (merged): https://github.com/code-corhuila/lms-access-api/pull/4
- ADR-007 gateway auth boundary: https://github.com/code-corhuila/library-docs/pull/8
- ADR-009 worker scheduling model: https://github.com/code-corhuila/library-docs/pull/10
- Commit history: https://github.com/code-corhuila/lms-access-api/commits/main/?author=Hermesss29
