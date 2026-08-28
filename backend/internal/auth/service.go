package auth

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidEmail       = errors.New("invalid email")
)

type User struct {
	ID    string
	Email string
	Role  string
}

type Service struct {
	db     *pgxpool.Pool
	tokens *TokenService
}

func NewService(db *pgxpool.Pool, tokens *TokenService) *Service {
	return &Service{db: db, tokens: tokens}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && strings.Contains(email, "@")
}

func (s *Service) Register(ctx context.Context, email, password, fullName, phone string) (User, error) {
	email = normalizeEmail(email)
	if !validateEmail(email) {
		return User{}, ErrInvalidEmail
	}
	if len(password) < 8 {
		return User{}, ErrWeakPassword
	}
	if strings.TrimSpace(fullName) == "" {
		return User{}, errors.New("full name is required")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	var user User
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'farmer')
		RETURNING id, email, role`, email, hash).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return User{}, ErrEmailExists
		}
		return User{}, err
	}

	_, err = s.db.Exec(ctx, `INSERT INTO farmer_profiles (user_id, full_name, phone) VALUES ($1, $2, NULLIF($3, ''))`, user.ID, strings.TrimSpace(fullName), strings.TrimSpace(phone))
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (User, string, error) {
	email = normalizeEmail(email)
	var user User
	var hash string
	err := s.db.QueryRow(ctx, `SELECT id, email, role, password_hash FROM users WHERE email = $1 AND is_active = TRUE`, email).Scan(&user.ID, &user.Email, &user.Role, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrInvalidCredentials
	}
	if err != nil {
		return User{}, "", err
	}
	if !CheckPassword(hash, password) {
		return User{}, "", ErrInvalidCredentials
	}

	token, err := s.tokens.Issue(user.ID, user.Role)
	if err != nil {
		return User{}, "", err
	}
	return user, token, nil
}

func (s *Service) UserByID(ctx context.Context, userID string) (User, error) {
	var user User
	err := s.db.QueryRow(ctx, `SELECT id, email, role FROM users WHERE id = $1 AND is_active = TRUE`, userID).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
