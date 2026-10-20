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

func TestFarmExpenseHandlerRequiresAuthentication(t *testing.T) {
	handler := NewFarmExpenseHandler(nil)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/expenses", strings.NewReader("{}"))
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

func TestFarmExpenseHandlerRejectsMalformedCycleID(t *testing.T) {
	handler := NewFarmExpenseHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crop-cycles/not-a-uuid/expenses", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	handler.List(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestFarmExpenseHandlerRejectsUnknownAndTrailingJSON(t *testing.T) {
	handler := NewFarmExpenseHandler(nil)
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	for _, body := range []string{
		`{"expense_type":"labor","description":"Harvest help","amount":100,"expense_date":"2026-10-19","extra":true}`,
		`{"expense_type":"labor","description":"Harvest help","amount":100,"expense_date":"2026-10-19"} {}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/expenses", strings.NewReader(body))
		req = req.WithContext(auth.WithUserID(context.Background(), userID))
		rec := httptest.NewRecorder()
		handler.Create(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, want %d", body, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestFarmExpenseHandlerCreateReturnsCreated(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	expenseID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_expenses.*WHERE cc.id = [$]1 AND fp.user_id = [$]2 RETURNING`).
		WithArgs(cycleID, userID, "labor", "Harvest help", 250000.0, "2026-10-19", "").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "expense_type", "description", "amount", "expense_date", "notes"}).
			AddRow(expenseID, cycleID, "labor", "Harvest help", 250000.0, "2026-10-19", ""))

	handler := NewFarmExpenseHandler(NewFarmExpenseService(pool))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/"+cycleID+"/expenses", strings.NewReader(`{"expense_type":"labor","description":"Harvest help","amount":250000,"expense_date":"2026-10-19"}`))
	req = req.WithContext(auth.WithUserID(context.Background(), userID))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), expenseID) {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
