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

```
cd lms-membership-api
copy <fork>\13-week\workshop\04-security-ci\authz_negative_test.go internal\adapter\in\httpapi\middleware\
go test ./internal/adapter/in/httpapi/middleware/ -run TestServiceToken -v
```

Expected today: **FAIL**. `RequireAuth` accepts any valid HS256 token for the whole `/students` group,
so the three requests reach the handler with 200. That red test is the open vulnerability of story #41;
it turns green once service tokens are limited to `GET /students/{id}` and `POST /students/{id}/suspend`.
Real console output is added in CP8.
