package farm

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestPlotServiceCreateScopesThroughAuthenticatedUser(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewPlotService(pool)
	userID := "10000000-0000-0000-0000-000000000001"
	farmID := "30000000-0000-0000-0000-000000000003"
	plotID := "40000000-0000-0000-0000-000000000004"
	pool.ExpectQuery(`INSERT INTO farm_plots \(farm_id, name, area_m2\)\s+SELECT f.id, \$3, \$4\s+FROM farms f\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE f.id = \$1 AND fp.user_id = \$2`).
		WithArgs(farmID, userID, "Blok A", 1000.0).
		WillReturnRows(pgxmock.NewRows([]string{"id", "farm_id", "name", "area_m2"}).AddRow(plotID, farmID, "Blok A", 1000.0))
	got, err := s.Create(context.Background(), userID, farmID, CreatePlotInput{Name: "Blok A", AreaM2: 1000})
	if err != nil { t.Fatal(err) }
	if got.ID != plotID || got.FarmID != farmID { t.Fatalf("unexpected plot: %+v", got) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPlotServiceRejectsInvalidArea(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewPlotService(pool)
	if _, err := s.Create(context.Background(), "user", "farm", CreatePlotInput{Name: "A", AreaM2: 0}); err == nil { t.Fatal("expected area validation error") }
}
