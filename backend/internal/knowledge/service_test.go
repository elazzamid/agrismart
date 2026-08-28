package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestPublishRequiresValidatedApprovedVersion(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewService(pool)
	documentID := "00000000-0000-0000-0000-000000000001"
	versionID := "00000000-0000-0000-0000-000000000002"
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1 AND d.status = 'validated' AND EXISTS`).
		WithArgs(documentID, versionID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	if err := s.Publish(context.Background(), documentID, versionID); err != nil { t.Fatal(err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPublishRejectsUnvalidatedDocument(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewService(pool)
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published'`).
		WithArgs("doc", "version").WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), "doc", "version"); err != ErrInvalidTransition { t.Fatalf("expected invalid transition, got %v", err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
