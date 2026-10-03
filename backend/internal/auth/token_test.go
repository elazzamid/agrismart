package auth

import (
	"os"
	"testing"
)

func TestTokenLifecycle(t *testing.T) {
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")

	service, err := NewTokenService()
	if err != nil {
		t.Fatalf("NewTokenService() error = %v", err)
	}

	token, err := service.Issue("user-123", "farmer")
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	userID, role, err := service.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if userID != "user-123" || role != "farmer" {
		t.Fatalf("unexpected claims: userID=%q role=%q", userID, role)
	}
}

func TestTokenServiceRejectsShortSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "short")
	defer os.Unsetenv("JWT_SECRET")

	if _, err := NewTokenService(); err == nil {
		t.Fatal("expected short JWT secret to be rejected")
	}
}
