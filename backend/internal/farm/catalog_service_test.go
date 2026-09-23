package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestCatalogServiceListCrops(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatalf("pgxmock.NewPool() error = %v", err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	pool.ExpectQuery(`SELECT id, code, name, COALESCE\(description, ''\), is_active FROM crops WHERE is_active = TRUE ORDER BY name`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "code", "name", "description", "is_active"}).
			AddRow("00000000-0000-0000-0000-000000000001", "CHILI", "Cabai", "Cabai merah", true))

	got, err := s.ListCrops(context.Background())
	if err != nil { t.Fatalf("ListCrops() error = %v", err) }
	if len(got) != 1 || got[0].Code != "CHILI" { t.Fatalf("unexpected crops: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatalf("expectations: %v", err) }
}

func TestCatalogServiceListVarietiesScopesCrop(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatalf("pgxmock.NewPool() error = %v", err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	cropID := "00000000-0000-0000-0000-000000000001"
	pool.ExpectQuery(`SELECT id, crop_id, name, COALESCE\(description, ''\), is_active FROM crop_varieties WHERE crop_id = \$1 AND is_active = TRUE ORDER BY name`).
		WithArgs(cropID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_id", "name", "description", "is_active"}).
			AddRow("00000000-0000-0000-0000-000000000002", cropID, "Varietas A", "", true))

	got, err := s.ListVarieties(context.Background(), cropID)
	if err != nil { t.Fatalf("ListVarieties() error = %v", err) }
	if len(got) != 1 || got[0].CropID != cropID { t.Fatalf("unexpected varieties: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatalf("expectations: %v", err) }
}

func TestCatalogServiceListGrowthStagesOrdersBySequence(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatalf("pgxmock.NewPool() error = %v", err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	cropID := "00000000-0000-0000-0000-000000000001"
	pool.ExpectQuery(`SELECT id, crop_id, name, sequence_no, min_days, max_days, COALESCE\(description, ''\) FROM crop_growth_stages WHERE crop_id = \$1 ORDER BY sequence_no`).
		WithArgs(cropID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_id", "name", "sequence_no", "min_days", "max_days", "description"}).
			AddRow("00000000-0000-0000-0000-000000000003", cropID, "Vegetatif", 1, 0, 30, "Awal pertumbuhan"))

	got, err := s.ListGrowthStages(context.Background(), cropID)
	if err != nil { t.Fatalf("ListGrowthStages() error = %v", err) }
	if len(got) != 1 || got[0].SequenceNo != 1 { t.Fatalf("unexpected stages: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatalf("expectations: %v", err) }
}
