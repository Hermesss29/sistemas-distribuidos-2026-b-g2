<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       05-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 05

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2
- SPRINT_GOAL: Define the system architecture (05-architecture) and the data model (06-data) for LMS-LIBRARY-V1 — ADRs, C4 views, and the PostgreSQL schema per bounded context.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-ARCH-001 | Architecture overview: C4 (System + Container), service catalog, principles, patterns, tech debt | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/blob/main/05-week/hu-status/05-architecture/overview.md |
| HU-ARCH-002 | ADR-001 documentation language, ADR-002 hexagonal modular monolith, ADR-003 NGINX reverse proxy | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/05-week/hu-status/05-architecture/decisions/records |
| HU-ARCH-003 | Cross-cutting concerns, deployment topology, pattern guide, STRIDE threat model | doing | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/05-week/hu-status/05-architecture |
| HU-DATA-001 | Data model per bounded context (administrators, students, books, loans) + ER diagram | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/blob/main/05-week/hu-status/06-data/models.md |
| HU-DATA-002 | Data dictionary, modeling conventions, migration strategy (golang-migrate) | done | https://github.com/Hermesss29/sistemas-distribuidos-2026-b-g2/tree/main/05-week/hu-status/06-data |
| HU-MVP-001 | Runnable MVP `lms-library-1.0.0` — Go services (access/catalog/membership), React SPA, NGINX, docker-compose, migrations, unit tests | done | `05-week/hu-status/MVP/lms-library-1.0.0.zip` |

## MVP delivered this week

A first end-to-end **MVP was built and packaged** as `05-week/hu-status/MVP/lms-library-1.0.0.zip`
(`lms-library` v1.0.0), turning the architecture and data docs into running code:

- **Backend (Go, Hexagonal / Ports & Adapters):** `access-service` (HU-01 login, JWT issuing,
  auth/correlation-id/CORS middleware, health endpoints), `catalog-service` (HU-04 create book,
  availability adjustment), `membership-service` (HU-02/HU-03 student create + management, with a
  circulation client port). Domain, application/usecase, and infrastructure layers separated;
  domain has no I/O.
- **Persistence:** PostgreSQL with `golang-migrate` migration files per service
  (`000001_create_*_table`, seed default administrator), matching `06-data/models.md`.
- **Frontend:** React + Vite + TypeScript + Tailwind SPA — login, protected routes, dashboard,
  books list/form, students list/form, loans list, overdue loans.
- **Infrastructure:** `docker-compose.yml` (services + Postgres + NGINX), `infra/nginx/nginx.conf`
  reverse proxy / gateway, `.env.example`, `SETUP.md`. Runs with `docker compose up -d --build`.
- **Tests:** unit tests on domain and usecases (`administrator_test.go`, `login_test.go`,
  `book_test.go`, `create_book_test.go`, `student_test.go`, `create_student_test.go`,
  `student_management_test.go`).
- **Deviation noted:** the MVP started extracting each bounded context into its own service
  (true microservices) rather than the single `library-api` modular monolith described in
  `05-architecture/overview.md` / `ADR-002`. This is a migration-in-progress (one domain per PR
  behind the NGINX gateway); the docs need an ADR update or amendment to reflect it.

## 2. My individual contribution
- Contributed to the MVP as **DevOps**: `docker-compose.yml`, NGINX reverse-proxy config
  (`infra/nginx/nginx.conf`), `.env.example` / environment-variable wiring, per-service
  `Dockerfile`/`Makefile`, and `SETUP.md`.
- Reviewed and validated the three ADRs (ADR-001/002/003) as part of the development team; confirmed the hexagonal modular monolith style and the single shared PostgreSQL database for v1, with the microservices-ready seams (module = bounded context).
- Contributed to the `06-data/models.md` schema: `administrators`, `students`, `books`, `loans` tables, their indexes, and the modeling decisions (soft delete on `students`, counter columns instead of a `Copy` entity, `ON DELETE RESTRICT` on loan FKs, `was_late` stored at return time).
- Cross-checked the data dictionary and ER diagram against `02-domain/entities-and-rules.md` invariants (INV-001 due_date = loan_date + 7d, INV-002 loan atomicity, suspension rules).
- Aligned architecture tech debt items AT-001 / AT-002 with the planned v2 evolution table.

## 3. Blockers and risks
- **MVP vs docs mismatch:** the MVP splits contexts into separate services while `ADR-002` still says single modular monolith. Needs a superseding/amending ADR before more services are extracted.
- MVP is partial: `circulation` (loans/returns/suspensions — the Core domain, HU-06/07/08) not yet extracted; loans UI is read-only.
- `07-api` OpenAPI contracts not started yet; several docs reference `07-api/contracts/openapi/` files that do not exist.
- `category` on `books` is free text (no controlled vocabulary) — may cause inconsistent catalog search (HU-05) until a fixed list is agreed.
- With per-service databases now in play, cross-service consistency (student ↔ loan) is no longer ACID — needs a defined strategy (event/outbox) not yet designed.

## 4. Plan for next week
- Write the ADR that reconciles the microservices split in the MVP with `ADR-002`.
- Extract the `circulation` service (loan registration HU-06, return HU-07, overdue HU-08) and wire it behind the NGINX gateway.
- Finish HU-ARCH-003: `cross-cutting.md`, `deployment.md`, `pattern-guide.md`, `security-threat-model.md`.
- Start the `07-api` OpenAPI contracts per service (API-First principle P1).
- Add integration tests (service + Postgres) on top of the current unit tests.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (`feat/HU-XX` -> `dev` -> `QA` -> `main`)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit) — integration tests still pending
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables (`.env.example`, no secrets committed)

## 6. Evidence links
- MVP package: 05-week/hu-status/MVP/lms-library-1.0.0.zip
- Architecture overview: 05-week/hu-status/05-architecture/overview.md
- ADRs: 05-week/hu-status/05-architecture/decisions/records/
- Hexagonal architecture: 05-week/hu-status/05-architecture/hexagonal-architecture.md
- Data model + ER diagram: 05-week/hu-status/06-data/models.md
- Data dictionary: 05-week/hu-status/06-data/data-dictionary.md
- Migration strategy: 05-week/hu-status/06-data/migration-strategy.md
