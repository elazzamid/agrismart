# Database

AgriSmart uses PostgreSQL as the transactional system of record.

## Migrations

SQL migrations are ordered and immutable after merge:

```text
migrations/
├── 000001_extensions.sql
├── 000002_identity.sql
├── 000003_farms.sql
├── 000004_crops.sql
├── 000005_crop_cycle_integrity.sql
├── 000006_knowledge_sources.sql
├── 000007_agricultural_knowledge.sql
├── 000008_opt_knowledge_links.sql
├── 000009_diagnosis_evidence.sql
├── 000010_fertilizer_nutrient_engine.sql
└── 000011_fertilizer_calculation_rules.sql
```

The migrations are intentionally plain SQL so they remain portable and easy to audit.

## Clean database verification

A clean PostgreSQL instance must accept all migrations in lexical order without manual schema changes.

## Seed data

Reference/seed data for the MVP crops is kept separate from schema migrations.

Current seed sets include:

- Core MVP crops and varieties.
- Cabai knowledge sources and source-backed knowledge fixtures.
- Nutrient reference data.
- MVP crop knowledge sources for padi and jagung.
- Evidence mappings and other knowledge fixtures.

Do not mix agricultural recommendations into core schema migration files. Actionable agricultural and pesticide guidance must remain source-backed and pass the knowledge validation workflow before publication.
