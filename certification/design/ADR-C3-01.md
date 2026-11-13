# ADR-C3-01: Replicate loan_db as a DDS replica set across three availability zones

## Status

Proposed

## Context

`loan_db` is the MongoDB database of `lms-circulation-api` (`ADR-005`). It holds the `loans` and
`idempotency_keys` collections, so it is the only place that knows which student has which book.
Section 4 of the design note picks it as the critical data store. Today it lives in the single MongoDB
instance that `lms-infra-mongo` runs on the same Docker host as everything else, and NFR-008 only
promises a daily backup:
RPO of 24 hours and RTO of 30 minutes for a database failure. If that host or its disk is lost, a
whole day of loans and returns disappears, and the copies stay reserved in `catalog_db` with no loan
behind them. The library desk needs an RPO close to zero and an RTO of a few minutes.

## Decision

I move `loan_db` to a DDS replica set in one Huawei Cloud region with one member in each of three
availability zones (primary in AZ1, secondaries in AZ2 and AZ3). `lms-circulation-api` connects with
the replica set URI and writes with write concern `majority`, so a loan is confirmed only when it is
stored in at least two AZs. If the primary or its AZ fails, the two remaining members elect a new
primary automatically. The daily automated backup to OBS stays as a second line against logical
errors, such as a bad delete that replication would copy to every member. Target: RPO 0 minutes,
RTO 5 minutes.

## Alternatives considered

- Keep the single MongoDB container and take daily backups (what NFR-008 describes today): rejected because a 24-hour RPO means losing a full day of loans, and restoring by hand takes longer than the desk can wait.
- DDS replica set across only two AZs: rejected because with two members a single failure leaves no majority, so no new primary is elected and writes stop until someone intervenes.
- Cross-region DR with a second region on warm standby: rejected for v1 because it roughly doubles the database cost for one university library, and an AZ failure is the realistic risk at this scale.

## Consequences

- Positive: an AZ failure loses no confirmed loan (RPO 0) and the service recovers by itself in about a minute, with no restore procedure.
- Positive: the code does not change much, because the CI and the compose file of `lms-circulation-db` already connect with `replicaSet=rs0`.
- Negative: three DDS nodes cost about three times one single-node instance, and `majority` writes add a few milliseconds of cross-AZ latency to `POST /loans`.
- Negative: the backups to OBS are still needed, so the team has to run and test restores anyway.
