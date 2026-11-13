# ADR-C3-02: Give each calling service its own internal token secret, stored in CSMS

## Status

Proposed

## Context

Section 2 of the design note (Security row) points to a gap the team already knows. `lms-membership-api`,
`lms-catalog-api`, `lms-circulation-api`, `lms-worker` and `lms-workflow` all sign or validate their
internal HS256 tokens with one shared `INTERNAL_JWT_SECRET`; the `.env.example` of each repository
declares the same variable. `ADR-013` added a per-route allow-list based on the
token `sub`, but with one shared secret any service can sign a token with any `sub`, so the
allow-list can be bypassed by anyone who has the secret. `ADR-013` left this as SEC-10, deferred
"until the system runs on more than one host". The target architecture runs on CCE nodes in two
AZs, so that condition is now met. Today the secret also lives in `.env` files on the host.

## Decision

Each calling service gets its own HS256 secret: one for `lms-circulation-api` (`sub=circulation-service`),
one for `lms-worker`, one for `lms-workflow`, and one more for each caller that appears later. The callee keeps a map from `sub` to secret and
validates the signature with the secret that belongs to the `sub` in the token, so a service can
only sign tokens as itself. The secrets are stored in CSMS, injected into each pod as environment
variables at deploy time (never in the image or in Git), and rotated every 90 days or at once if
one leaks. During a rotation the callee accepts the old and the new secret for one hour.

## Alternatives considered

- Keep the shared secret and rely only on the ADR-013 allow-list: rejected because the allow-list trusts the `sub`, and with a shared secret the `sub` can be forged.
- Move internal calls to RS256 with a key pair per service: rejected for now because every Go service would need key distribution and JWKS lookups, which is more work than the team can do this term. It stays as the next step if the number of callers grows.
- Mutual TLS between services with a service mesh on CCE: rejected because it adds a mesh to operate, and the team's principle P5 avoids infrastructure it cannot run.

## Consequences

- Positive: a leaked secret only exposes the routes that one caller is allowed to use, and the ADR-013 allow-list can finally be trusted.
- Positive: removing the secrets from `.env` files closes one of the leak paths named in the week-13 threat model.
- Negative: the callee has to manage a map of secrets, and every new caller means one more secret to create, inject and rotate.
- Negative: a rotation done in the wrong order can break service-to-service calls, so it needs the one-hour overlap and a short runbook.
