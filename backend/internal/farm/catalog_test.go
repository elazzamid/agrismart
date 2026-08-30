package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestCatalogServiceListCropsReturnsActiveCrops(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	pool.ExpectQuery(`SELECT id, code, name, COALESCE\(description, ''\)\s+FROM crops WHERE is_active = TRUE ORDER BY name`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "code", "name", "description"}).
			AddRow("crop-1", "padi", "Padi", "Tanaman padi"))

	items, err := s.ListCrops(context.Background())
	if err != nil { t.Fatal(err) }
	if len(items) != 1 || items[0].Code != "padi" { t.Fatalf("unexpected crops: %+v", items) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestCatalogServiceListVarietiesScopesToCrop(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	cropID := "00000000-0000-0000-0000-000000000001"
	pool.ExpectQuery(`SELECT id, crop_id, name, COALESCE\(description, ''\)\s+FROM crop_varieties\s+WHERE crop_id = \$1 AND is_active = TRUE ORDER BY name`).
		WithArgs(cropID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_id", "name", "description"}).
			AddRow("variety-1", cropID, "IR64", "Varietas padi"))

	items, err := s.ListVarieties(context.Background(), cropID)
	if err != nil { t.Fatal(err) }
	if len(items) != 1 || items[0].CropID != cropID { t.Fatalf("unexpected varieties: %+v", items) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestCatalogServiceListGrowthStagesOrdersBySequence(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()

	s := NewCatalogService(pool)
	cropID := "00000000-0000-0000-0000-000000000001"
	pool.ExpectQuery(`SELECT id, crop_id, name, sequence_no, min_days, max_days, COALESCE\(description, ''\)\s+FROM crop_growth_stages\s+WHERE crop_id = \$1 ORDER BY sequence_no`).
		WithArgs(cropID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_id", "name", "sequence_no", "min_days", "max_days", "description"}).
			AddRow("stage-1", cropID, "Persemaian", 1, 0, 20, "Tahap awal"))

	items, err := s.ListGrowthStages(context.Background(), cropID)
	if err != nil { t.Fatal(err) }
	if len(items) != 1 || items[0].SequenceNo != 1 { t.Fatalf("unexpected stages: %+v", items) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
