## Summary

Adds `07-api/guidelines.md`, a file this folder's README lists as required and that didn't exist. It sets the REST rules that sit above the OpenAPI contracts, then audits the five existing contracts against those rules.

Eleven decisions closed (D-G01–D-G11): versioning, pagination parameters and envelope, date ranges, error body shape, error-code catalogue, 403 in a single-role system, the 400/409/422 boundary, malformed UUIDs, the missing 500 declarations, and the two 401 codes. Each states the decision, the reason and the consequences.

## Defects found in §3

| # | Defect | Owner |
|---|---|---|
| DF-1 | Seven of eight operations declare no 500; only `POST /auth/login` does | Each service owner |
| DF-2 | `GET /students/{id}` declares no 400, though it's the only endpoint taking a UUID in the path | membership-service |
| DF-3 | The gateway's 503 references `ErrorResponse` but names no `error` code | api-gateway |
| DF-4 | `POST /books` still lives in `library-api.yaml` while `catalog-service.yaml` doesn't exist | catalog-service |
| DF-5 | Handlers fall back to `limit = 20` above the maximum instead of clamping to 100 | Service owners (surfaced by #17) |

§4 tracks six gaps (H-1–H-6), each with an owner. H-6 is mine.

## Relationship to PR #17

They are **not disjoint**: both create `07-api/guidelines.md`. This is individual work, so we audited the folder separately — #17 from the implemented handlers in the `lms-*-api` repos, this one from the OpenAPI contracts and the services not yet built.

**Where they agree:** the pagination envelope. #17 adopts D-G03 from this PR as the forward decision.

**Where they diverged:** #17's audit showed the handlers fall back to `limit = 20` rather than clamping to 100. This document keeps the clamp, matching D-C5 in the course reference, and records the difference as DF-5.

**Authority over the file:** whichever PR merges first owns `07-api/guidelines.md`; the second rebases and keeps only its own delta. If #17 lands first, that's mine to do.

## Test plan

- [x] Checked against the five existing contracts: `_shared.yaml`, `access-service.yaml`, `membership-service.yaml`, `api-gateway.yaml`, `library-api.yaml`
- [x] Decisions not yet observable are listed as gaps in §4, not presented as current behaviour
- [x] All internal anchors verified
- [x] No conflicting guidance with `docs/07-api-guidelines` — reconciled with #17 as described above; the remaining overlap is the file path, resolved by rebase at merge time
- [x] 395 lines, within the 400-line cap. No further commits planned on this branch
