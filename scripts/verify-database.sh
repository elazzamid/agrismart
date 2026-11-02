#!/usr/bin/env sh
set -eu

if [ -z "${DB_CONTAINER:-}" ]; then
  DB_CONTAINER="$(docker compose ps -q postgres | sed -n '1p')"
fi

[ -n "$DB_CONTAINER" ] || {
  echo "PostgreSQL Compose service is not running; start it with 'docker compose up -d postgres' or set DB_CONTAINER." >&2
  exit 1
}
DB_USER="${POSTGRES_USER:-agrismart}"
DB_NAME="${POSTGRES_DB:-agrismart}"

printf '%s\n' "Checking PostgreSQL container: $DB_CONTAINER"
docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME"

printf '%s\n' "Checking required tables..."
TABLE_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name IN ('users','farmer_profiles','farms','farm_plots','crops','crop_varieties','crop_growth_stages','crop_cycles','farm_activities','farm_inputs','farm_expenses','farm_harvests');")"

[ "$TABLE_COUNT" = "12" ] || {
  echo "Expected 12 required tables, found $TABLE_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking seeded MVP crop catalog..."
MVP_CROP_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(DISTINCT code) FROM crops WHERE is_active = TRUE AND code IN ('cabai','padi','jagung');")"

[ "$MVP_CROP_COUNT" = "3" ] || {
  echo "Expected 3 active MVP crops (cabai, padi, jagung), found $MVP_CROP_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking seeded MVP crop varieties..."
MVP_VARIETY_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM (VALUES ('padi','IR64'),('padi','Inpari 32'),('jagung','Bisi 18'),('jagung','NK 212'),('cabai','Cabai Merah Keriting'),('cabai','Cabai Rawit')) AS expected(code,name) JOIN crops c ON c.code = expected.code AND c.is_active = TRUE JOIN crop_varieties v ON v.crop_id = c.id AND v.name = expected.name AND v.is_active = TRUE;")"

[ "$MVP_VARIETY_COUNT" = "6" ] || {
  echo "Expected 6 seeded MVP crop varieties, found $MVP_VARIETY_COUNT" >&2
  exit 1
}

printf '%s\n' "Checking seeded MVP crop growth stages..."
MVP_GROWTH_STAGE_COUNT="$(docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM (VALUES ('padi','Persemaian',1),('padi','Vegetatif',2),('padi','Generatif',3),('padi','Pematangan',4),('jagung','Perkecambahan',1),('jagung','Vegetatif',2),('jagung','Pembungaan',3),('jagung','Pengisian Biji',4),('cabai','Persemaian',1),('cabai','Vegetatif',2),('cabai','Berbunga',3),('cabai','Pembuahan',4),('cabai','Panen',5)) AS expected(code,name,sequence_no) JOIN crops c ON c.code = expected.code AND c.is_active = TRUE JOIN crop_growth_stages s ON s.crop_id = c.id AND s.name = expected.name AND s.sequence_no = expected.sequence_no;")"

[ "$MVP_GROWTH_STAGE_COUNT" = "13" ] || {
  echo "Expected 13 seeded MVP crop growth stages, found $MVP_GROWTH_STAGE_COUNT" >&2
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
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_harvests_quantity_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_harvests_quantity_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_harvests_unit_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_harvests_unit_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_harvests_quality_grade_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_harvests_quality_grade_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'farm_harvests_notes_check') THEN
    RAISE EXCEPTION 'Missing required constraint: farm_harvests_notes_check';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND indexname = 'idx_farm_harvests_cycle_date') THEN
    RAISE EXCEPTION 'Missing required index: idx_farm_harvests_cycle_date';
  END IF;
END
$verify$;
SQL

printf '%s\n' "Database verification: PASS"
