package farm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
	"github.com/pashagolub/pgxmock/v4"
)

func TestCreateCropCycleRequiresAuthentication(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/plot-1/crop-cycles", strings.NewReader("{\"crop_id\":\"crop-1\",\"planting_date\":\"2026-08-01\"}"))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestCreateCropCycleRejectsInvalidDate(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/plot-1/crop-cycles", strings.NewReader("{\"crop_id\":\"crop-1\",\"planting_date\":\"2026-02-30\"}"))
	req = req.WithContext(auth.WithUserID(context.Background(), "farmer-1"))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestCreateCropCycleReturnsCreatedForOwnedPlot(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	farmerID := "00000000-0000-0000-0000-000000000001"
	plotID := "00000000-0000-0000-0000-000000000002"
	cropID := "00000000-0000-0000-0000-000000000003"
	cycleID := "00000000-0000-0000-0000-000000000004"
	pool.ExpectQuery(`INSERT INTO crop_cycles \(plot_id, crop_id, variety_id, planting_date\)\s+SELECT p.id, \$3, \$4, \$5::date\s+FROM farm_plots p JOIN farms f ON f.id = p.farm_id\s+WHERE p.id = \$1 AND f.farmer_id = \$2`).
		WithArgs(plotID, farmerID, cropID, (*string)(nil), "2026-08-01").
		WillReturnRows(pgxmock.NewRows([]string{"id", "plot_id", "crop_id", "variety_id", "planting_date", "status"}).
			AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	body := "{\"crop_id\":\"" + cropID + "\",\"planting_date\":\"2026-08-01\"}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/"+plotID+"/crop-cycles", strings.NewReader(body))
	req = req.WithContext(auth.WithUserID(context.Background(), farmerID))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), cycleID) {
		t.Fatalf("response does not contain crop cycle id: %s", rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
