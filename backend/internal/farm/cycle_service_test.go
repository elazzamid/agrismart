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
        WithArgs(plotID, farmerID, cropID, nil, "2026-08-01").
        WillReturnRows(pgxmock.NewRows([]string{"id","plot_id","crop_id","variety_id","planting_date","status"}).AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))
    got, err := s.Create(context.Background(), farmerID, plotID, CreateCropCycleInput{CropID: cropID, PlantingDate: "2026-08-01"})
    if err != nil { t.Fatal(err) }
    if got.ID != cycleID || got.PlotID != plotID { t.Fatalf("unexpected cycle: %+v", got) }
    if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
