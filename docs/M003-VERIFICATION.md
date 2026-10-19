# M003.1 Farm Activities Verification

## Scope

This slice adds activity records to an existing crop cycle. It does not implement farm inputs, expenses, or harvest records.

## API

All endpoints require a valid Bearer token.

- `GET /api/v1/crop-cycles/{id}/activities` — list activity records for a crop cycle owned by the authenticated farmer.
- `POST /api/v1/crop-cycles/{id}/activities` — create an activity record for a crop cycle owned by the authenticated farmer.

Example request:

```json
{
  "activity_type": "fertilizing",
  "activity_date": "2026-10-19",
  "notes": "Pemupukan susulan"
}
```

Validation rules:

- `activity_type` is required and limited to 80 characters.
- `activity_date` is required and must use the valid `YYYY-MM-DD` format.
- `notes` is optional and limited to 2,000 characters.
- Unknown JSON fields are rejected.
- A crop cycle that does not belong to the authenticated farmer returns HTTP 404.

## Required verification

- Service tests verify ownership-scoped inserts and list access.
- Handler tests verify authentication, malformed cycle IDs, creation, and hiding unowned cycles.
- `go test ./... -count=1` and `go build ./cmd/api` must pass in CI.
- The database smoke test must confirm the new migration creates `farm_activities` and its constraints.
