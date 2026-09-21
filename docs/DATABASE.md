# AgriSmart Initial Database Model

## M001.1 purpose

This is the domain model baseline. It is intentionally smaller than the final production schema. Tables are introduced only when their ownership and relationships are understood.

## Identity

### users
- id UUID PK
- email TEXT UNIQUE NOT NULL
- password_hash TEXT NOT NULL
- role TEXT NOT NULL
- is_active BOOLEAN NOT NULL DEFAULT true
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

### farmer_profiles
- id UUID PK
- user_id UUID UNIQUE FK users(id)
- full_name TEXT NOT NULL
- phone TEXT
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

## Farm

### farms
- id UUID PK
- farmer_id UUID FK farmer_profiles(id)
- name TEXT NOT NULL
- location_name TEXT
- latitude NUMERIC(9,6)
- longitude NUMERIC(9,6)
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

### farm_plots
- id UUID PK
- farm_id UUID FK farms(id)
- name TEXT NOT NULL
- area_m2 NUMERIC(14,2) NOT NULL
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

### crop_cycles
- id UUID PK
- plot_id UUID FK farm_plots(id)
- crop_id UUID FK crops(id)
- variety_id UUID NULL FK crop_varieties(id)
- planting_date DATE NOT NULL
- status TEXT NOT NULL
- created_at TIMESTAMPTZ NOT NULL
- updated_at TIMESTAMPTZ NOT NULL

## Knowledge

### crops
- id UUID PK
- code TEXT UNIQUE NOT NULL
- name TEXT NOT NULL
- description TEXT
- is_active BOOLEAN NOT NULL DEFAULT true

### crop_varieties
- id UUID PK
- crop_id UUID FK crops(id)
- name TEXT NOT NULL
- description TEXT
- is_active BOOLEAN NOT NULL DEFAULT true

### crop_growth_stages
- id UUID PK
- crop_id UUID FK crops(id)
- name TEXT NOT NULL
- sequence_no INT NOT NULL
- min_days INT
- max_days INT
- description TEXT

Future fertilizer, pesticide, pest, disease, guide, and knowledge-document tables will be added after the core schema is migrated and verified.

## Rules

1. Use UUID identifiers for domain entities.
2. Use UTC timestamps.
3. Foreign keys are mandatory for owned relationships.
4. Soft deletion is not introduced globally; add it only where business requirements justify it.
5. Migration files are immutable after merge.
6. Seed data is separate from schema migrations when practical.
7. Agricultural recommendations must be versionable and traceable to sources.
