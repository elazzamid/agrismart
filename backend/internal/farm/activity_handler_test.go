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

func TestFarmActivityHandlerRequiresAuthentication(t *testing.T) {
	handler := NewActivityHandler(nil)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/activities", strings.NewReader(`{"activity_type":"watering","activity_date":"2026-10-19"}`))
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

func TestFarmActivityHandlerRejectsMalformedCropCycleID(t *testing.T) {
	handler := NewActivityHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crop-cycles/not-a-uuid/activities", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestFarmActivityHandlerRejectsTrailingJSON(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	handler := NewActivityHandler(NewFarmActivityService(pool))
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"

	for _, body := range []string{
		`{"activity_type":"watering","activity_date":"2026-10-19"} {"unexpected":true}`,
		`{"activity_type":"watering","activity_date":"2026-10-19"} trailing`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/activities", strings.NewReader(body))
		req = req.WithContext(auth.WithUserID(context.Background(), userID))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, want %d; response=%s", body, rec.Code, http.StatusBadRequest, rec.Body.String())
		}
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFarmActivityHandlerCreateReturnsCreated(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	activityID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_activities.*WHERE cc.id = \\$1 AND fp.user_id = \\$2 RETURNING`).
		WithArgs(cycleID, userID, "watering", "2026-10-19", "Morning").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "activity_type", "activity_date", "notes"}).
			AddRow(activityID, cycleID, "watering", "2026-10-19", "Morning"))

	handler := NewActivityHandler(NewFarmActivityService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/activities", strings.NewReader(`{"activity_type":"watering","activity_date":"2026-10-19","notes":"Morning"}`))
	req = req.WithContext(auth.WithUserID(context.Background(), userID))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), activityID) {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFarmActivityHandlerListReturnsNotFoundForUnownedCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	cycleID := "00000000-0000-0000-0000-000000000002"
	pool.ExpectQuery(`SELECT EXISTS.*WHERE cc.id = \\$1 AND fp.user_id = \\$2`).
		WithArgs(cycleID, "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	handler := NewActivityHandler(NewFarmActivityService(pool))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crop-cycles/"+cycleID+"/activities", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
