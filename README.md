# hello-api

Developer-owned application source. A small Go HTTP service (standard library only).

| Endpoint | Purpose |
|---|---|
| `GET /healthz` | liveness/readiness probe |
| `GET /` | JSON: service name and `ENVIRONMENT` |
| `GET /metrics` | Prometheus counter `http_requests_total{status="2xx".."5xx"}` |

Repository: https://github.com/__GH_USER__/hello-api
Image (from Phase 2 CI): `ghcr.io/__GH_USER_LC__/hello-api`

Run the tests the same way CI will, in the build container:

    docker build -t hello-api:dev .
