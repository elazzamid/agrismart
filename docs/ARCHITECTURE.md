# AgriSmart Architecture Specification

## Status
M001.1 — Technical Foundation

## Decisions
- Backend: Go modular monolith
- Database: PostgreSQL
- API: REST/JSON over HTTPS
- Frontend: Next.js, mobile-first web
- Local development: Docker Compose
- Database changes: versioned SQL migrations
- Authentication: secure session/JWT-ready foundation; implementation must keep auth isolated from domain modules
- File storage: object-storage abstraction; local filesystem only for development
- AI: separate application module, not part of M001

## Repository

```text
agrismart/
├── backend/
├── frontend/
├── database/
│   └── migrations/
├── docs/
├── infra/
├── scripts/
├── .github/
│   └── workflows/
├── docker-compose.yml
├── .env.example
└── README.md
```

## Backend modules

```text
backend/
├── cmd/api/
├── internal/
│   ├── auth/
│   ├── farmer/
│   ├── farm/
│   ├── crop/
│   ├── knowledge/
│   ├── fertilizer/
│   ├── pesticide/
│   ├── pest/
│   ├── disease/
│   ├── activity/
│   ├── finance/
│   └── platform/
└── migrations/
```

Business modules must not import HTTP handlers directly. Transport, application/service logic, and persistence concerns remain separated enough to permit testing without a running web server.

## Initial domain boundaries

### Identity
- users
- farmer_profiles

### Farm
- farms
- farm_plots
- crop_cycles

### Agricultural knowledge
- crops
- crop_varieties
- crop_growth_stages
- fertilizers
- fertilizer_nutrients
- pesticides
- active_ingredients
- pests
- diseases
- weeds
- cultivation_guides
- management_guides

### Operations
- farm_activities
- farm_inputs
- farm_expenses
- harvests

AI, weather, marketplace, and advanced diagnosis remain outside M001.

## API baseline

Base path:

`/api/v1`

Health endpoint:

`GET /health`

Expected deterministic response:

```json
{"status":"ok"}
```

API errors should use a consistent JSON envelope and never expose internal stack traces or secrets.

## Configuration

Environment variables are documented in `.env.example`. Secrets must never be committed. Production values must be supplied by deployment infrastructure.

## Database

PostgreSQL is the system of record for transactional data. UUID primary keys are preferred for externally exposed domain identifiers. Timestamps use UTC. Foreign keys and database constraints enforce invariants where practical.

Migrations are forward-only and numbered. A clean database must be reproducible from migrations alone.

## Testing

M001 minimum verification:
- backend unit tests
- API health test
- authentication foundation tests
- migration smoke test
- frontend build/typecheck
- static checks

## Deployment principle

Development must be reproducible through Docker Compose. Production deployment may use containers, but production-specific infrastructure is not part of M001.
