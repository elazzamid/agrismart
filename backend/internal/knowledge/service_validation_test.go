package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestValidateApprovedSetsValidatedStatus(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer db.Close()
	s := NewService(db)
	ctx := context.Background()
	db.ExpectBegin()
	db.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM knowledge_versions`).WithArgs("v1", "d1").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	db.ExpectExec(`INSERT INTO knowledge_validations`).WithArgs("v1", "expert", "approved", "ok").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	db.ExpectExec(`UPDATE knowledge_documents SET status = \$2`).WithArgs("d1", "validated").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	db.ExpectCommit()
	if err := s.Validate(ctx, "d1", "v1", "expert", "approved", "ok"); err != nil { t.Fatal(err) }
	if err := db.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPublishRequiresValidatedDocumentAndApprovedVersion(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer db.Close()
	s := NewService(db)
	db.ExpectExec(`UPDATE knowledge_documents d SET status = 'published'`).WithArgs("d1", "v1").WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), "d1", "v1"); err != ErrInvalidTransition { t.Fatalf("expected ErrInvalidTransition, got %v", err) }
	if err := db.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPublishRequiresSourceBackedLatestApprovedVersion(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer db.Close()
	s := NewService(db)
	db.ExpectExec(`(?s)^UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d\.id = \$1 AND d\.status = 'validated' AND EXISTS \(.*v\.source_id IS NOT NULL.*\)$`).
		WithArgs("d1", "v1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	if err := s.Publish(context.Background(), "d1", "v1"); err != nil { t.Fatal(err) }
	if err := db.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
