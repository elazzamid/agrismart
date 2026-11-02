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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/00000000-0000-0000-0000-000000000002/crop-cycles", strings.NewReader("{\"crop_id\":\"00000000-0000-0000-0000-000000000003\",\"planting_date\":\"2026-08-01\"}"))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestCreateCropCycleRejectsUnavailableVarietyAsBadRequest(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	userID := "00000000-0000-0000-0000-000000000001"
	cropID := "00000000-0000-0000-0000-000000000003"
	varietyID := "00000000-0000-0000-0000-000000000005"
	pool.ExpectQuery(validateCropSelectionQuery).
		WithArgs(cropID, &varietyID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	body := "{\"crop_id\":\"" + cropID + "\",\"variety_id\":\"" + varietyID + "\",\"planting_date\":\"2026-08-01\"}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/00000000-0000-0000-0000-000000000002/crop-cycles", strings.NewReader(body))
	req = req.WithContext(auth.WithUserID(context.Background(), userID))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "variety_id must be active") {
		t.Fatalf("response does not explain invalid variety: %s", rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCropCycleRejectsInvalidDate(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/00000000-0000-0000-0000-000000000002/crop-cycles", strings.NewReader("{\"crop_id\":\"00000000-0000-0000-0000-000000000003\",\"planting_date\":\"2026-02-30\"}"))
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCropCycleReturnsCreatedForOwnedPlot(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	userID := "00000000-0000-0000-0000-000000000001"
	plotID := "00000000-0000-0000-0000-000000000002"
	cropID := "00000000-0000-0000-0000-000000000003"
	cycleID := "00000000-0000-0000-0000-000000000004"
	pool.ExpectQuery(validateCropSelectionQuery).
		WithArgs(cropID, (*string)(nil)).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	pool.ExpectQuery(`INSERT INTO crop_cycles \(plot_id, crop_id, variety_id, planting_date\)\s+SELECT p.id, \$3, \$4, \$5::date\s+FROM farm_plots p\s+JOIN farms f ON f.id = p.farm_id\s+JOIN farmer_profiles fp ON fp.id = f.farmer_id\s+WHERE p.id = \$1 AND fp.user_id = \$2\s+RETURNING id, plot_id, crop_id, variety_id, planting_date::text, status`).
		WithArgs(plotID, userID, cropID, (*string)(nil), "2026-08-01").
		WillReturnRows(pgxmock.NewRows([]string{"id", "plot_id", "crop_id", "variety_id", "planting_date", "status"}).
			AddRow(cycleID, plotID, cropID, nil, "2026-08-01", "active"))
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	body := "{\"crop_id\":\"" + cropID + "\",\"planting_date\":\"2026-08-01\"}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/"+plotID+"/crop-cycles", strings.NewReader(body))
	req = req.WithContext(auth.WithUserID(context.Background(), userID))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), cycleID) {
		t.Fatalf("response does not contain crop cycle id: %s", rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}


func TestListCropCyclesRejectsMalformedPlotIDBeforeDatabaseCall(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/plots/not-a-uuid/crop-cycles", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	h.ListCropCycles(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "plot id must be a valid UUID") {
		t.Fatalf("response does not explain invalid plot id: %s", rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCropCycleRejectsMalformedPlotIDBeforeDatabaseCall(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	h := NewHandler(nil, nil, nil, NewCropCycleService(pool))
	body := "{\"crop_id\":\"00000000-0000-0000-0000-000000000003\",\"planting_date\":\"2026-08-01\"}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/plots/not-a-uuid/crop-cycles", strings.NewReader(body))
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	h.CreateCropCycle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "plot id must be a valid UUID") {
		t.Fatalf("response does not explain invalid plot id: %s", rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
