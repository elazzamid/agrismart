package farm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestFarmInputHandlerRequiresAuthentication(t *testing.T) {
	handler := NewFarmInputHandler(nil)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/inputs", strings.NewReader("{}"))
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

func TestFarmInputHandlerRejectsMalformedCycleID(t *testing.T) {
	handler := NewFarmInputHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crop-cycles/not-a-uuid/inputs", nil)
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()

	handler.List(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestFarmInputHandlerRejectsUnknownAndTrailingJSON(t *testing.T) {
	handler := NewFarmInputHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/inputs", strings.NewReader(`{"input_type":"fertilizer","product_name":"Urea","quantity":1,"unit":"kg","applied_date":"2026-10-19","extra":true}`))
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/crop-cycles/00000000-0000-0000-0000-000000000001/inputs", strings.NewReader(`{"input_type":"fertilizer","product_name":"Urea","quantity":1,"unit":"kg","applied_date":"2026-10-19"} {}`))
	req = req.WithContext(auth.WithUserID(context.Background(), "user-1"))
	rec = httptest.NewRecorder()
	handler.Create(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("trailing JSON status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
