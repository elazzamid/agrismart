package farm

import (
    "context"
    "testing"

    "github.com/pashagolub/pgxmock/v4"
)

func TestCropCycleCreateScopesThroughPlotOwner(t *testing.T) {
    pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }
    defer pool.Close()
    s := NewCropCycleService(pool)
    farmerID := "00000000-0000-0000-0000-000000000001"
    plotID := "00000000-0000-0000-0000-000000000002"
    cropID := "00000000-0000-0000-0000-000000000003"
    cycleID := "00000000-0000-0000-0000-000000000004"
    pool.ExpectQuery(`INSERT INTO crop_cycles \(plot_id, crop_id, variety_id, planting_date\)\s+SELECT p.id, \$3, \$4, \$5::date\s+FROM farm_plots p JOIN farms f ON f.id = p.farm_id\s+WHERE p.id = \$1 AND f.farmer_id = \$2`).
        WithArgs(plotID, farmerID, cropID, (*string)(nil), "2026-08-01").
        WillReturnRows(pgxmock.NewRows([]string{"id","plot_id","crop_id","variety_id","planting_date","status"}).AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))
    got, err := s.Create(context.Background(), farmerID, plotID, CreateCropCycleInput{CropID: cropID, PlantingDate: "2026-08-01"})
    if err != nil { t.Fatal(err) }
    if got.ID != cycleID || got.PlotID != plotID { t.Fatalf("unexpected cycle: %+v", got) }
    if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}


func TestCropCycleCreateRejectsInvalidInputBeforeDatabaseCall(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewCropCycleService(pool)

	for _, in := range []CreateCropCycleInput{
		{PlantingDate: "2026-08-01"},
		{CropID: "00000000-0000-0000-0000-000000000003", PlantingDate: "2026-02-30"},
	} {
		if _, err := s.Create(context.Background(), "farmer", "plot", in); err != ErrInvalidCropCycleInput {
			t.Fatalf("Create(%+v) error = %v, want ErrInvalidCropCycleInput", in, err)
		}
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestCropCycleListScopesThroughPlotOwner(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewCropCycleService(pool)
	farmerID := "00000000-0000-0000-0000-000000000001"
	plotID := "00000000-0000-0000-0000-000000000002"
	cycleID := "00000000-0000-0000-0000-000000000004"
	pool.ExpectQuery(`SELECT cc.id, cc.plot_id, cc.crop_id, cc.variety_id, cc.planting_date::text, cc.status\s+FROM crop_cycles cc JOIN farm_plots p ON p.id = cc.plot_id\s+JOIN farms f ON f.id = p.farm_id\s+WHERE cc.plot_id = \$1 AND f.farmer_id = \$2 ORDER BY cc.planting_date DESC`).
		WithArgs(plotID, farmerID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "plot_id", "crop_id", "variety_id", "planting_date", "status"}).
			AddRow(cycleID, plotID, "00000000-0000-0000-0000-000000000003", nil, "2026-08-01", "active"))

	got, err := s.List(context.Background(), farmerID, plotID)
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].ID != cycleID { t.Fatalf("unexpected crop cycles: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
