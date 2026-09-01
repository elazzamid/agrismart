#!/usr/bin/env sh
set -eu

DB_USER="${POSTGRES_USER:-agrismart}"
DB_NAME="${POSTGRES_DB:-agrismart}"

psql_exec() {
  docker compose exec -T postgres psql -U "$DB_USER" -d "$DB_NAME" "$@"
}

pg_isready() {
  docker compose exec -T postgres pg_isready -U "$DB_USER" -d "$DB_NAME"
}

printf '%s\n' "Checking PostgreSQL service: postgres"
pg_isready

printf '%s\n' "Checking required tables..."
TABLE_COUNT="$(psql_exec -Atc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name IN ('users','farmer_profiles','farms','farm_plots','crops','crop_varieties','crop_growth_stages','crop_cycles');")"

[ "$TABLE_COUNT" = "8" ] || {
  echo "Expected 8 core tables, found $TABLE_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking core and knowledge constraints..."
CONSTRAINT_COUNT="$(psql_exec -Atc "SELECT count(*) FROM pg_constraint WHERE conname IN ('users_role_check','farms_latitude_check','farms_longitude_check','farm_plots_area_check','crop_varieties_crop_name_unique','crop_growth_stages_sequence_unique','crop_cycles_status_check','fertilizer_nutrients_percentage_check','fertilizer_nutrients_percentage_range_check');")"

[ "$CONSTRAINT_COUNT" = "9" ] || {
  echo "Expected 9 required constraints, found $CONSTRAINT_COUNT" >&2
  exit 1
}

printf '%s\n' "Database verification: PASS"
