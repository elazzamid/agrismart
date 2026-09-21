# Docker Development

## PostgreSQL

The development database is provided by Docker Compose using PostgreSQL 16.

Start it with:

```bash
docker compose up -d postgres
```

The SQL files in `database/migrations/` are mounted into PostgreSQL's initialization directory and are executed in lexical filename order on a fresh database volume.

## Verification

Run:

```bash
sh scripts/verify-database.sh
```

The verification checks PostgreSQL readiness, the eight core M001 tables, and the key database constraints.

To force a clean migration run locally:

```bash
docker compose down -v
docker compose up -d postgres
sh scripts/verify-database.sh
```

Never commit a real `.env` file or production credentials. Use `.env.example` as the local configuration template.
