# Database

AgriSmart uses PostgreSQL as the transactional system of record.

## Migrations

SQL migrations are ordered and immutable after merge:

```text
migrations/
├── 000001_extensions.sql
├── 000002_identity.sql
├── 000003_farms.sql
└── 000004_crops.sql
```

The migrations are intentionally plain SQL so they remain portable and easy to audit.

## Clean database verification

A clean PostgreSQL instance must accept the migrations in lexical order without manual schema changes.

## Seed data

Reference/seed data for the MVP crops will be added separately after the schema is verified. Do not mix agricultural recommendations into the core schema migration files.
