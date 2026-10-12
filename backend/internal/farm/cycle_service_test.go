package farm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
	"github.com/pashagolub/pgxmock/v4"
)

func TestCropCycleCreateScopesThroughAuthenticatedUserProfile(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewCropCycleService(pool)
	userID := "10000000-0000-0000-0000-000000000001"
	profileID := "20000000-0000-0000-0000-000000000002"
	plotID := "30000000-0000-0000-0000-000000000003"
	cropID := "40000000-0000-0000-0000-000000000004"
	cycleID := "50000000-0000-0000-0000-000000000005"

	pool.ExpectQuery(`INSERT INTO crop_cycles \(plot_id, crop_id, variety_id, planting_date\)\s+SELECT p.id, \$3, \$4, \$5::date\s+FROM farm_plots p\s+JOIN farms f ON f.id = p.farm_id\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE p.id = \$1 AND fp.user_id = \$2\s+RETURNING id, plot_id, crop_id, variety_id, planting_date::text, status`).
		WithArgs(plotID, userID, cropID, (*string)(nil), "2026-08-01").
		WillReturnRows(pgxmock.NewRows([]string{"id", "plot_id", "crop_id", "variety_id", "planting_date", "status"}).
			AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))

	got, err := s.Create(context.Background(), userID, plotID, CreateCropCycleInput{CropID: cropID, PlantingDate: "2026-08-01"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != cycleID || got.PlotID != plotID || userID == profileID {
		t.Fatalf("unexpected crop cycle or test fixture: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}


func TestCreateCropCycleHandlerRequiresAuthentication(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/plot-id/crop-cycles", strings.NewReader(`{"crop_id":"crop-id","planting_date":"2026-08-01"}`))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestCreateCropCycleHandlerUsesAuthenticatedUserAndPlotPath(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	userID := "10000000-0000-0000-0000-000000000001"
	plotID := "30000000-0000-0000-0000-000000000003"
	cropID := "40000000-0000-0000-0000-000000000004"
	cycleID := "50000000-0000-0000-0000-000000000005"
	pool.ExpectQuery(`INSERT INTO crop_cycles \(plot_id, crop_id, variety_id, planting_date\)\s+SELECT p.id, \$3, \$4, \$5::date\s+FROM farm_plots p\s+JOIN farms f ON f.id = p.farm_id\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE p.id = \$1 AND fp.user_id = \$2\s+RETURNING id, plot_id, crop_id, variety_id, planting_date::text, status`).
		WithArgs(plotID, userID, cropID, (*string)(nil), "2026-08-01").
		WillReturnRows(pgxmock.NewRows([]string{"id", "plot_id", "crop_id", "variety_id", "planting_date", "status"}).
			AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))

	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/"+plotID+"/crop-cycles", strings.NewReader(`{"crop_id":"`+cropID+`","planting_date":"2026-08-01"}`))
	req = req.WithContext(auth.WithUserID(req.Context(), userID))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var got CropCycle
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ID != cycleID || got.PlotID != plotID || got.CropID != cropID {
		t.Fatalf("unexpected crop cycle response: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
