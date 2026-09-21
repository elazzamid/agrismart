package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "correct horse battery staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == password {
		t.Fatal("password must not be stored as plaintext")
	}
	if !CheckPassword(hash, password) {
		t.Fatal("CheckPassword() rejected the original password")
	}
	if CheckPassword(hash, "wrong password") {
		t.Fatal("CheckPassword() accepted an incorrect password")
	}
}

func TestHashPasswordRejectsEmptyPassword(t *testing.T) {
	if _, err := HashPassword(""); err != ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
