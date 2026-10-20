# M003 Farm Operations Verification

## Scope

This branch implements three farm-operation slices:

- **M003.1 — Farm Activities:** record and list activities for a crop cycle.
- **M003.2 — Farm Inputs:** record and list inputs actually applied to a crop cycle. This is historical recordkeeping, not an agronomic recommendation or dose calculator.
- **M003.3 — Farm Expenses:** record and list crop-cycle operating costs. This is expense history, not financial advice.

## API

All endpoints require a valid Bearer token and enforce that the crop cycle belongs to the authenticated farmer.

### Farm activities

- `GET /api/v1/crop-cycles/{id}/activities`
- `POST /api/v1/crop-cycles/{id}/activities`

Request fields: `activity_type`, `activity_date` (`YYYY-MM-DD`), and optional `notes`.

### Farm inputs

- `GET /api/v1/crop-cycles/{id}/inputs`
- `POST /api/v1/crop-cycles/{id}/inputs`

Request fields: `input_type` (`fertilizer`, `pesticide`, `seed`, `other`), `product_name`, positive `quantity`, `unit`, `applied_date` (`YYYY-MM-DD`), and optional `notes`.

### Farm expenses

- `GET /api/v1/crop-cycles/{id}/expenses`
- `POST /api/v1/crop-cycles/{id}/expenses`

Example request:

```json
{
  "expense_type": "labor",
  "description": "Harvest help",
  "amount": 250000,
  "expense_date": "2026-10-19",
  "notes": "Paid in cash"
}
```

Allowed expense types: `seed`, `fertilizer`, `pesticide`, `labor`, `irrigation`, `equipment`, `transport`, `other`. Description is 1–160 characters; amount must be between Rp0.01 and Rp999,999,999,999.99; date must be a valid `YYYY-MM-DD` date; notes are limited to 2,000 characters.

## Validation and access rules

- Crop-cycle IDs must be UUIDs.
- Unknown JSON fields and trailing JSON values are rejected.
- Missing or unowned crop cycles return HTTP 404.
- Inputs and expenses are historical farmer records; the API does not calculate agricultural doses or provide financial recommendations.

## Database migration order

M002 knowledge migrations use `000006` through `000011`. Farm operations use `000012_farm_activities.sql`, `000013_farm_inputs.sql`, and `000014_farm_expenses.sql`.

## Required verification

- Backend tests: `cd backend && go test ./... -count=1`
- Backend build: `cd backend && go build ./cmd/api`
- Database migration smoke test and `scripts/verify-database.sh` must pass in CI.
