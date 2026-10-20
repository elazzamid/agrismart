# Database

AgriSmart uses PostgreSQL as the transactional system of record.

## Migration order

Migrations are plain SQL and must apply in lexical order on a clean database. Migration files are immutable after merge.

| Migration range | Scope |
| --- | --- |
| `000001`–`000005` | Core schema, identity, farms, crops, and crop-cycle integrity |
| `000006`–`000011` | M002 agricultural knowledge foundation |
| `000012` | M003.1 farm activity records |
| `000013` | M003.2 farm input records |

The M003 migration numbers intentionally follow M002. Do not reuse a migration number across features or renumber a migration after it has been merged.

## Clean database verification

A clean PostgreSQL instance must accept all migrations in lexical order without manual schema changes. CI starts PostgreSQL, applies the SQL files, and runs `scripts/verify-database.sh` to check required tables, constraints, and indexes.

## Seed data

Reference and knowledge seed data is kept separate from schema migrations when practical. Do not mix agricultural recommendations into core schema migration files.
