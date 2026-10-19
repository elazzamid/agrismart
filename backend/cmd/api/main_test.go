package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubDatabasePinger struct {
	err error
}

func (s stubDatabasePinger) Ping(context.Context) error {
	return s.err
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json content type, got %q", got)
	}
	if got := rec.Body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected response body: %q", got)
	}
}

func TestReadinessHandlerReturnsReadyWhenDatabaseResponds(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ready", nil)
	rec := httptest.NewRecorder()

	readinessHandler(rec, req, stubDatabasePinger{})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "{\"status\":\"ready\"}\n" {
		t.Fatalf("unexpected response body: %q", got)
	}
}

func TestReadinessHandlerReturnsServiceUnavailableWhenDatabaseFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ready", nil)
	rec := httptest.NewRecorder()

	readinessHandler(rec, req, stubDatabasePinger{err: errors.New("database unavailable")})

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if got := rec.Body.String(); got != "{\"status\":\"not_ready\"}\n" {
		t.Fatalf("readiness response must not expose internal error, got %q", got)
	}
}
