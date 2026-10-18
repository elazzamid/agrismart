package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func expectLockedStatus(pool pgxmock.PgxPoolIface, status string) {
	pool.ExpectBegin()
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1 FOR UPDATE`).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow(status))
}

func expectPublishUpdate(pool pgxmock.PgxPoolIface, result int64) {
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published'`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", result))
}

func TestPublishRequiresValidatedApprovedVersion(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil { t.Fatal(err) }
	defer pool.Close()
	s := NewService(pool)
	documentID, versionID := "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"
	pool.ExpectBegin()
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1 FOR UPDATE`).WithArgs(documentID).
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("validated"))
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1`).
		WithArgs(documentID, versionID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()
	if err := s.Publish(context.Background(), documentID, versionID); err != nil { t.Fatal(err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPublishRejectsUnvalidatedDocument(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }; defer pool.Close()
	s := NewService(pool)
	expectLockedStatus(pool, "review")
	pool.ExpectRollback()
	if err := s.Publish(context.Background(), "doc", "version"); err != ErrInvalidTransition { t.Fatalf("expected invalid transition, got %v", err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestPublishRequiresLatestVersion(t *testing.T) {
 expectPublishRejected(t, `(?s).*NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL.*v.id = \$2.*NULLIF\(BTRIM\(v.source_id\), ''\) IS NOT NULL.*kv.decision = 'approved'.*NOT EXISTS.*v.version_no = \(SELECT MAX\(v2.version_no\).*`)
}
func TestPublishRejectsVersionWithNewerNonApprovedValidation(t *testing.T) {
 expectPublishRejected(t, `(?s).*kv.decision = 'approved'.*NOT EXISTS \(\s*SELECT 1 FROM knowledge_validations kv2\s*WHERE kv2.version_id = v.id AND \(kv2.validated_at, kv2.id\) > \(kv.validated_at, kv.id\)\s*\).*`)
}
func TestPublishRejectsMissingAuthor(t *testing.T) {
 expectPublishRejected(t, `(?s).*WHERE d.id = \$1 AND d.status = 'validated' AND NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL AND EXISTS.*`)
}
func TestPublishRejectsBlankProvenanceSource(t *testing.T) {
 expectPublishRejected(t, `(?s).*v.id = \$2 AND v.document_id = d.id AND NULLIF\(BTRIM\(v.source_id\), ''\) IS NOT NULL AND kv.decision = 'approved'.*`)
}

func expectPublishRejected(t *testing.T, sqlPattern string) {
	t.Helper()
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }; defer pool.Close()
	s := NewService(pool)
	expectLockedStatus(pool, "validated")
	pool.ExpectExec(sqlPattern).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	pool.ExpectRollback()
	if err := s.Publish(context.Background(), "doc", "version"); err != ErrInvalidTransition { t.Fatalf("expected invalid transition, got %v", err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestAddVersionRejectsValidatedDocument(t *testing.T) { expectAddVersionRejected(t, "validated") }
func TestAddVersionRejectsPublishedDocument(t *testing.T) { expectAddVersionRejected(t, "published") }

func expectAddVersionRejected(t *testing.T, status string) {
	t.Helper()
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }; defer pool.Close()
	s := NewService(pool)
	pool.ExpectBegin()
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1 FOR UPDATE`).WithArgs("doc").
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow(status))
	pool.ExpectRollback()
	if err := s.AddVersion(context.Background(), "doc", 2, "content", nil); err != ErrInvalidTransition { t.Fatalf("expected invalid transition, got %v", err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

func TestValidateApprovedOlderVersionKeepsDocumentInReview(t *testing.T) {
	pool, err := pgxmock.NewPool(); if err != nil { t.Fatal(err) }; defer pool.Close()
	s := NewService(pool)
	pool.ExpectBegin()
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1 FOR UPDATE`).WithArgs("doc").
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("review"))
	pool.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM knowledge_versions WHERE id = \$1 AND document_id = \$2\)`).
		WithArgs("version-1", "doc").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	pool.ExpectExec(`INSERT INTO knowledge_validations`).WithArgs("version-1", "validator", "approved", "").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	pool.ExpectQuery(`SELECT EXISTS\(\s*SELECT 1 FROM knowledge_versions v`).WithArgs("version-1", "doc").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
	pool.ExpectExec(`UPDATE knowledge_documents SET status = \$2, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("doc", "review").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()
	if err := s.Validate(context.Background(), "doc", "version-1", "validator", "approved", ""); err != nil { t.Fatal(err) }
	if err := pool.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}
