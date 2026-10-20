# M003 Farm Operations Verification

## Scope

This branch implements the first two farm-operation slices:

- **M003.1 — Farm Activities:** record and list activities for a crop cycle.
- **M003.2 — Farm Inputs:** record and list inputs actually applied to a crop cycle. This is historical recordkeeping, not an agronomic recommendation or dose calculator.

## API

All endpoints require a valid Bearer token and enforce that the crop cycle belongs to the authenticated farmer.

### Farm activities

- `GET /api/v1/crop-cycles/{id}/activities`
- `POST /api/v1/crop-cycles/{id}/activities`

Example request:

```json
{
  "activity_type": "fertilizing",
  "activity_date": "2026-10-19",
  "notes": "Pemupukan susulan"
}
```

### Farm inputs

- `GET /api/v1/crop-cycles/{id}/inputs`
- `POST /api/v1/crop-cycles/{id}/inputs`

Example request:

```json
{
  "input_type": "fertilizer",
  "product_name": "Urea",
  "quantity": 1.25,
  "unit": "kg",
  "applied_date": "2026-10-19",
  "notes": "Pemupukan susulan"
}
```

## Validation and access rules

- Crop-cycle IDs must be UUIDs.
- Unknown JSON fields and trailing JSON values are rejected.
- Activity type is required and limited to 80 characters; activity date must be a valid `YYYY-MM-DD` date; notes are limited to 2,000 characters.
- Input type must be one of `fertilizer`, `pesticide`, `seed`, or `other`.
- Input product name is 1–120 characters; quantity must be finite and at least 0.001; unit is 1–20 characters; applied date must be a valid `YYYY-MM-DD` date; notes are limited to 2,000 characters.
- Missing or unowned crop cycles return HTTP 404.

## Database migration order

M002 knowledge migrations use `000006` through `000011`. To avoid a filename collision with M002's `000006_knowledge_sources.sql`, farm activities use `000012_farm_activities.sql` and farm inputs use `000013_farm_inputs.sql`.

## Required verification

- Backend tests: `cd backend && go test ./... -count=1`
- Backend build: `cd backend && go build ./cmd/api`
- Database migration smoke test and `scripts/verify-database.sh` must pass in CI.
