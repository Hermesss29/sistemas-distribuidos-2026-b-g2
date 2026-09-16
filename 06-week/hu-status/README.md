<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       06-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 06

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Define the API contracts (07-api) for LMS-LIBRARY-V1 following the API-First principle, and present the Cut 1 progress in the class review session.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-API-001 | Shared OpenAPI components (pagination, error format, correlation ID) | done | [07-api/_shared.yaml](./07-api/_shared.yaml) |
| HU-API-002 | Access service API contract (administrator login, JWT issuing, health) | done | [07-api/access-service.yaml](./07-api/access-service.yaml) |
| HU-API-003 | Membership service API contract (student registration and management) | done | [07-api/membership-service.yaml](./07-api/membership-service.yaml) |

## Sprint review — Cut 1 progress presentation

The team presented the Cut 1 progress in class on Monday, 7 September 2026, during the 8:40 a.m. session. My participation:

- **Jira board walkthrough:** explained how the team's board is organized, how work items are registered and tracked, and how the activities were split among the three members so each one owns a defined scope and the work can advance in parallel.
- **Pull Request workflow:** explained that every change goes through a Pull Request reviewed by the team before being integrated, as the mechanism the team uses for change control and for traceability of who contributed what.
- **Peer feedback:** gave feedback to classmates on their own projects during the session.

## 2. My individual contribution
- Contributed with the team to defining the `07-api` OpenAPI contracts, closing the gap reported as a blocker in week 05 ("07-api OpenAPI contracts not started yet").
- `_shared.yaml`: shared component library referenced by every service contract — `PaginatedMeta` for list responses and a project-specific `ErrorResponse` aligned with the actual `response.ErrorBody` returned by the services, including the correlation ID propagated to the structured logs.
- `access-service.yaml`: contract for the Access bounded context — administrator authentication, JWT issuing and health endpoints, documented against the API gateway base URL.
- `membership-service.yaml`: contract for the Membership bounded context — student registration, listing and per-student operations, plus the service health endpoint.
- Documented the deliberate scope exclusions in the access contract (no registration endpoint, no refresh token or logout, no JWKS, no roles in the token), each justified against the v1 operating model of a single Administrator.
- Presented the team's Jira board and the Pull Request workflow during the Cut 1 review session, and gave feedback to classmates on their projects.

## 3. Blockers and risks
- The contracts describe the intended API surface; they still need to be validated against the running services and kept in sync as the endpoints evolve.
- The `circulation` context (loans, returns, penalties) has no contract yet, since the service has not been extracted.

## 4. Plan for next week
- Write the contracts for the remaining services and align them with the API gateway routes.
- Validate the contracts against the running MVP.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Shared components: [06-week/hu-status/07-api/_shared.yaml](./07-api/_shared.yaml)
- Access service contract: [06-week/hu-status/07-api/access-service.yaml](./07-api/access-service.yaml)
- Membership service contract: [06-week/hu-status/07-api/membership-service.yaml](./07-api/membership-service.yaml)
