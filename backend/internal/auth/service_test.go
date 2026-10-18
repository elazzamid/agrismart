package auth

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsEmailUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "email unique constraint",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "users_email_key",
			},
			want: true,
		},
		{
			name: "different unique constraint",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "farmer_profiles_user_id_key",
			},
			want: false,
		},
		{
			name: "different postgres error",
			err: &pgconn.PgError{
				Code:           "23503",
				ConstraintName: "users_email_key",
			},
			want: false,
		},
		{
			name: "wrapped email unique constraint",
			err: errors.Join(errors.New("insert failed"), &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "users_email_key",
			}),
			want: true,
		},
		{
			name: "plain error mentioning duplicate key",
			err:  errors.New("duplicate key value violates constraint"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmailUniqueViolation(tt.err); got != tt.want {
				t.Fatalf("isEmailUniqueViolation() = %v, want %v", got, tt.want)
			}
		})
	}
}
