package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterMissingFullNameReturnsBadRequest(t *testing.T) {
	handler := NewHandler(NewService(nil, nil))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(`{"email":"farmer@example.com","password":"strong-pass-123","full_name":"  "}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for missing full_name, got %d; body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "full name is required") {
		t.Fatalf("expected validation error in response, got %s", rec.Body.String())
	}
}

func TestRegisterRejectsTrailingJSONValue(t *testing.T) {
	handler := NewHandler(NewService(nil, nil))
	body := `{"email":"farmer@example.com","password":"strong-pass-123","full_name":"Farmer"} {"unexpected":"second value"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for trailing JSON value, got %d; body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestLoginRejectsTrailingGarbage(t *testing.T) {
	handler := NewHandler(NewService(nil, nil))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"farmer@example.com","password":"strong-pass-123"} trailing`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for trailing garbage, got %d; body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}
