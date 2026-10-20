#!/usr/bin/env sh
set -eu

DB_CONTAINER="${DB_CONTAINER:-agrismart-postgres}"
DB_USER="${POSTGRES_USER:-agrismart}"
DB_NAME="${POSTGRES_DB:-agrismart}"

printf '%s\n' "Checking PostgreSQL container: $DB_CONTAINER"
docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME"

printf '%s\n' "Checking required tables..."
TABLE_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name IN ('users','farmer_profiles','farms','farm_plots','crops','crop_varieties','crop_growth_stages','crop_cycles','farm_activities','farm_inputs','farm_expenses');")"

[ "$TABLE_COUNT" = "11" ] || {
  echo "Expected 11 required tables, found $TABLE_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking required constraints and indexes..."
docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 <<'SQL'
DO $verify$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'users_role_check') THEN
    RAISE EXCEPTION 'Missing required constraint: users_role_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farms_latitude_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farms_latitude_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farms_longitude_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farms_longitude_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_plots_area_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_plots_area_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'crop_varieties_crop_name_unique') THEN
    RAISE EXCEPTION 'Missing required constraint: crop_varieties_crop_name_unique';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'crop_growth_stages_sequence_unique') THEN
    RAISE EXCEPTION 'Missing required constraint: crop_growth_stages_sequence_unique';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'crop_cycles_status_check') THEN
    RAISE EXCEPTION 'Missing required constraint: crop_cycles_status_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_activities_type_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_activities_type_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_activities_notes_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_activities_notes_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND indexname = 'idx_farm_activities_cycle_date') THEN
    RAISE EXCEPTION 'Missing required index: idx_farm_activities_cycle_date';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_inputs_type_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_inputs_type_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_inputs_product_name_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_inputs_product_name_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_inputs_quantity_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_inputs_quantity_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_inputs_unit_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_inputs_unit_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_inputs_notes_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_inputs_notes_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND indexname = 'idx_farm_inputs_cycle_date') THEN
    RAISE EXCEPTION 'Missing required index: idx_farm_inputs_cycle_date';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_expenses_type_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_expenses_type_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_expenses_description_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_expenses_description_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_expenses_amount_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_expenses_amount_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_expenses_notes_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_expenses_notes_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND indexname = 'idx_farm_expenses_cycle_date') THEN
    RAISE EXCEPTION 'Missing required index: idx_farm_expenses_cycle_date';
  END IF;
END
$verify$;
SQL

printf '%s\n' "Database verification: PASS"
