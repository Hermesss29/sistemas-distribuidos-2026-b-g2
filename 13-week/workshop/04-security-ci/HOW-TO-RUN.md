# How to run - security gates and negative authorization test (CP4)

| File | What it is |
|---|---|
| `security.yml` | GitHub Actions workflow for lms-membership-api: gitleaks, govulncheck and the negative authorization test, on every pull request to develop, qa and main |
| `authz_negative_test.go` | Go test for `RequireAuth` in `internal/adapter/in/httpapi/middleware/` |

## Story and criterion it covers

[library-docs#41](https://github.com/code-corhuila/library-docs/issues/41), STRIDE E on the boundary
circulation-api -> membership-api: "An HS256 service token calling `POST /api/v1/students/{id}/deactivate`
gets 403", plus `PATCH /api/v1/students/{id}` and `GET /api/v1/students`.

## Run it locally

The test runs against the real `RequireAuth` of lms-membership-api (branch `develop`). The file is copied
in only to run it and is removed afterwards: it is not committed to the code repo.

```
cd C:\Users\hjuni\OneDrive\Escritorio\ReposVarios\lms-qa\repos\lms-membership-api
git checkout develop
git pull
copy C:\Users\hjuni\OneDrive\Escritorio\Distribuidos\sistemas-distribuidos-2026-b-g2\13-week\workshop\04-security-ci\authz_negative_test.go internal\adapter\in\httpapi\middleware\
go test ./internal/adapter/in/httpapi/middleware/ -run TestServiceToken -v
del internal\adapter\in\httpapi\middleware\authz_negative_test.go
```

## Real output (2026-10-30, my machine)

```
PS C:\Users\hjuni\OneDrive\Escritorio\ReposVarios\lms-qa\repos\lms-membership-api> go test ./internal/adapter/in/httpapi/middleware/ -run TestServiceToken -v
=== RUN   TestServiceTokenCannotReachAdminRoutes
=== RUN   TestServiceTokenCannotReachAdminRoutes/POST_/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97/deactivate
    authz_negative_test.go:53: service token reached POST /api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97/deactivate: status 200, want 403
=== RUN   TestServiceTokenCannotReachAdminRoutes/PATCH_/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97
    authz_negative_test.go:53: service token reached PATCH /api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97: status 200, want 403
=== RUN   TestServiceTokenCannotReachAdminRoutes/GET_/api/v1/students
    authz_negative_test.go:53: service token reached GET /api/v1/students: status 200, want 403
--- FAIL: TestServiceTokenCannotReachAdminRoutes (0.00s)
    --- FAIL: TestServiceTokenCannotReachAdminRoutes/POST_/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97/deactivate (0.00s)
    --- FAIL: TestServiceTokenCannotReachAdminRoutes/PATCH_/api/v1/students/5f2c8e1a-7b3d-4c9e-8a2f-1d6b0e4c3a97 (0.00s)
    --- FAIL: TestServiceTokenCannotReachAdminRoutes/GET_/api/v1/students (0.00s)
FAIL
FAIL    github.com/code-corhuila/lms-membership-api/internal/adapter/in/httpapi/middleware      1.390s
FAIL
```

**Result: FAIL, as expected.** The service token minted like circulation-api's (`sub: circulation-service`,
HS256, 1 min `exp`) reaches `POST /students/{id}/deactivate`, `PATCH /students/{id}` and `GET /students`
with 200. That proves the open vulnerability of story #41 (STRIDE E). The test turns green when
`RequireAuth` limits service tokens to `GET /students/{id}` and `POST /students/{id}/suspend`, and from then
on the `authz` job of `security.yml` blocks any PR that opens those routes again.

## Workflow validation

`security.yml` checked with **actionlint 1.7.7**: no errors (exit code 0).
