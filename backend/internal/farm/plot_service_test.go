package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestPlotServiceCreateScopesThroughFarmOwner(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewPlotService(pool)
	farmID := "00000000-0000-0000-0000-000000000001"
	farmerID := "00000000-0000-0000-0000-000000000002"
	plotID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_plots \(farm_id, name, area_m2\)\s+SELECT id, \$3, \$4 FROM farms WHERE id = \$1 AND farmer_id = \$2`).
		WithArgs(farmID, farmerID, "Blok A", 1000.0).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farm_id", "name", "area_m2"}).AddRow(plotID, farmID, "Blok A", 1000.0))
	got, err := s.Create(context.Background(), farmerID, farmID, CreatePlotInput{Name: "Blok A", AreaM2: 1000})
	if err != nil { t.Fatal(err) }
	if got.ID != plotID || got.FarmID != farmID { t.Fatalf("unexpected plot: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPlotServiceRejectsInvalidArea(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewPlotService(pool)
	if _, err := s.Create(context.Background(), "farmer", "farm", CreatePlotInput{Name: "A", AreaM2: 0}); err == nil { t.Fatal("expected area validation error") }
}
