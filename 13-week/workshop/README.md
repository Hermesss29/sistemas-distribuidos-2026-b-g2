# Week 13 Workshop - Threat Modeling and Security Backlog

Applied to the LMS (Library Management System) of team G2 (code-corhuila).

<!-- CONFIG-START -->
- FULL_NAME: Hermes Pascuas Herrera
- GITHUB_USER: Hermesss29
- TEAM: G2 - LMS (code-corhuila)
- DOCS_REPO: https://github.com/code-corhuila/library-docs
- REFLECTION_DOC: https://docs.google.com/document/d/1lgmE3VAZHeHtZZbuf14te5gzRz_AdVDH8Z5tWrZ_7tU/edit?usp=sharing
- BOUNDARY: circulation-api -> membership-api (internal service-to-service REST over lms-network, HS256 service token)
- CLAIM_ISSUE: https://github.com/code-corhuila/library-docs/issues/33
- LIBRARY_BOOK: Pressman, R. S. Ingenieria del software: un enfoque practico, 7th ed., McGraw-Hill, 2010. Call number 005.3 P727A. Pages 318, 378, 454
- LIBRARY_RECORD_TIME: 2026-10-29 08:40
- STRIDE_PR: https://github.com/code-corhuila/library-docs/pull/40
- STORIES: https://github.com/code-corhuila/library-docs/issues/41 (E, risk 6, P1: service tokens limited to the routes circulation-api needs), https://github.com/code-corhuila/library-docs/issues/42 (S, risk 6, P1: verify sub, aud and exp of service tokens), https://github.com/code-corhuila/library-docs/issues/43 (I, risk 6, P1: eligibility-only response for circulation-api)
<!-- CONFIG-END -->

## Checkpoints

| CP | Where | Status |
|---|---|---|
| CP0 | README.md + claim issue in library-docs | done |
| CP1 | 01-attack-surface.md | done |
| CP2 | library-docs PR: 05-architecture/security/stride-Hermesss29.md | done |
| CP3 | library-docs issues + team Project | done |
| CP4 | 04-security-ci/ | done |
| CP6 | library-docs PR #40: review answered, STRIDE finished | done |
| CP7 | issues #41-#43 refined, STORIES line | done |
