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

func TestFarmHarvestHandlerRequiresAuthentication(t *testing.T) {
	handler := NewFarmHarvestHandler(nil)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/harvests", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		if method == http.MethodGet {
			handler.List(rec, req)
		} else {
			handler.Create(rec, req)
		}
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s status = %d, want %d", method, rec.Code, http.StatusUnauthorized)
		}
	}
}

func TestFarmHarvestHandlerRejectsMalformedCycleID(t *testing.T) {
	handler := NewFarmHarvestHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crop-cycles/not-a-uuid/harvests", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	handler.List(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestFarmHarvestHandlerRejectsUnknownAndTrailingJSON(t *testing.T) {
	handler := NewFarmHarvestHandler(nil)
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	for _, body := range []string{
		`{"quantity":125.5,"unit":"kg","harvest_date":"2026-10-19","unexpected":true}`,
		`{"quantity":125.5,"unit":"kg","harvest_date":"2026-10-19"} {}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/harvests", strings.NewReader(body))
		req = req.WithContext(auth.WithUserID(context.Background(), userID))
		rec := httptest.NewRecorder()
		handler.Create(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, want %d", body, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestFarmHarvestHandlerCreateReturnsCreated(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	harvestID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_harvests.*WHERE cc.id = [$]1 AND fp.user_id = [$]2 RETURNING`).
		WithArgs(cycleID, userID, 125.5, "kg", "2026-10-19", "Grade A", "First pick").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "quantity", "unit", "harvest_date", "quality_grade", "notes"}).
			AddRow(harvestID, cycleID, 125.5, "kg", "2026-10-19", "Grade A", "First pick"))

	handler := NewFarmHarvestHandler(NewFarmHarvestService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/harvests", strings.NewReader(`{"quantity":125.5,"unit":"kg","harvest_date":"2026-10-19","quality_grade":"Grade A","notes":"First pick"}`))
	req = req.WithContext(auth.WithUserID(context.Background(), userID))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), harvestID) {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
