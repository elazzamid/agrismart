#!/usr/bin/env sh
set -eu

DB_CONTAINER="${DB_CONTAINER:-agrismart-postgres}"
DB_USER="${POSTGRES_USER:-agrismart}"
DB_NAME="${POSTGRES_DB:-agrismart}"

printf '%s\n' "Checking PostgreSQL container: $DB_CONTAINER"
docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME"

printf '%s\n' "Checking required tables..."
TABLE_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name IN ('users','farmer_profiles','farms','farm_plots','crops','crop_varieties','crop_growth_stages','crop_cycles');")"

[ "$TABLE_COUNT" = "8" ] || {
  echo "Expected 8 core tables, found $TABLE_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking core constraints..."
docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 <<'SQL'
SELECT 1 FROM pg_constraint WHERE conname = 'users_role_check';
SELECT 1 FROM pg_constraint WHERE conname = 'farms_latitude_check';
SELECT 1 FROM pg_constraint WHERE conname = 'farms_longitude_check';
SELECT 1 FROM pg_constraint WHERE conname = 'farm_plots_area_check';
SELECT 1 FROM pg_constraint WHERE conname = 'crop_varieties_crop_name_unique';
SELECT 1 FROM pg_constraint WHERE conname = 'crop_growth_stages_sequence_unique';
SELECT 1 FROM pg_constraint WHERE conname = 'crop_cycles_status_check';

printf '%s\n' "Checking knowledge constraints..."
SELECT 1 FROM pg_constraint WHERE conname = 'fertilizer_nutrients_percentage_check';
SELECT 1 FROM pg_constraint WHERE conname = 'fertilizer_nutrients_percentage_range_check';
SQL

printf '%s\n' "Database verification: PASS"
