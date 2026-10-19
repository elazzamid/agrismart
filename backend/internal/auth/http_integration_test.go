package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// authIntegrationPool uses the CI PostgreSQL database when available. These
// tests are skipped for local unit-test runs without DATABASE_URL.
func authIntegrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set; skipping authentication integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping PostgreSQL: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func seedAuthMiddlewareUser(t *testing.T, pool *pgxpool.Pool, role string, active bool) string {
	t.Helper()
	ctx := context.Background()
	email := fmt.Sprintf("auth-middleware-%d@example.test", time.Now().UnixNano())
	var userID string
	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, role, is_active)
		VALUES ($1, 'unused-test-hash', $2, $3)
		RETURNING id::text
	`, email, role, active).Scan(&userID)
	if err != nil {
		t.Fatalf("insert test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID)
	})
	return userID
}

func TestAuthenticatedRejectsInactiveUser(t *testing.T) {
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	pool := authIntegrationPool(t)
	userID := seedAuthMiddlewareUser(t, pool, "farmer", false)

	tokens, err := NewTokenService()
	if err != nil {
		t.Fatalf("NewTokenService(): %v", err)
	}
	token, err := tokens.Issue(userID, "farmer")
	if err != nil {
		t.Fatalf("Issue(): %v", err)
	}

	handler := NewHandler(NewService(pool, tokens)).Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not run for an inactive user")
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}

func TestAuthenticatedUsesCurrentDatabaseRole(t *testing.T) {
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	pool := authIntegrationPool(t)
	// The token deliberately contains the old role. The active DB record is
	// authoritative and should supply the role passed to protected handlers.
	userID := seedAuthMiddlewareUser(t, pool, "admin", true)

	tokens, err := NewTokenService()
	if err != nil {
		t.Fatalf("NewTokenService(): %v", err)
	}
	token, err := tokens.Issue(userID, "farmer")
	if err != nil {
		t.Fatalf("Issue(): %v", err)
	}

	handler := NewHandler(NewService(pool, tokens)).Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID, ok := UserIDFromContext(r.Context())
		if !ok || gotID != userID {
			t.Errorf("context user ID = %q, present=%v; want %q", gotID, ok, userID)
		}
		gotRole, ok := RoleFromContext(r.Context())
		if !ok || gotRole != "admin" {
			t.Errorf("context role = %q, present=%v; want admin", gotRole, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}
