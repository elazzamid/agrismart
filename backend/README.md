# AgriSmart Backend

Go 1.22 modular-monolith foundation.

## Run locally

```bash
go test ./...
go run ./cmd/api
```

The server listens on `:8080` by default. Set `API_ADDR` to override it.

## Health endpoint

```text
GET /api/v1/health
```

Response:

```json
{"status":"ok"}
```

The M001 backend intentionally contains no agricultural business logic yet.
