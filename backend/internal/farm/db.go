package farm

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DBTX captures the pgx operations used by farm services and is implemented
// by both *pgxpool.Pool and pgxmock.PgxPoolIface.
type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
