package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestPublishRequiresValidatedApprovedVersion(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	documentID := "00000000-0000-0000-0000-000000000001"
	versionID := "00000000-0000-0000-0000-000000000002"
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1 AND d.status = 'validated' AND NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL AND EXISTS`).
		WithArgs(documentID, versionID).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	if err := s.Publish(context.Background(), documentID, versionID); err != nil {
		t.Fatal(err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishRejectsUnvalidatedDocument(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published'`).
		WithArgs("doc", "version").WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), "doc", "version"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishRequiresLatestVersion(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	documentID := "00000000-0000-0000-0000-000000000001"
	versionID := "00000000-0000-0000-0000-000000000002"
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1 AND d.status = 'validated' AND NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL AND EXISTS \(\s*SELECT 1 FROM knowledge_versions v\s*JOIN knowledge_validations kv ON kv.version_id = v.id\s*WHERE v.id = \$2 AND v.document_id = d.id AND v.source_id IS NOT NULL AND kv.decision = 'approved'\s*AND NOT EXISTS \(\s*SELECT 1 FROM knowledge_validations kv2\s*WHERE kv2.version_id = v.id AND \(kv2.validated_at, kv2.id\) > \(kv.validated_at, kv.id\)\s*\)\s*AND v.version_no = \(SELECT MAX\(v2.version_no\) FROM knowledge_versions v2 WHERE v2.document_id = d.id\)`).
		WithArgs(documentID, versionID).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), documentID, versionID); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition for non-latest version, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishRejectsVersionWithNewerNonApprovedValidation(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	documentID := "00000000-0000-0000-0000-000000000001"
	versionID := "00000000-0000-0000-0000-000000000002"
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1 AND d.status = 'validated' AND NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL AND EXISTS \(\s*SELECT 1 FROM knowledge_versions v\s*JOIN knowledge_validations kv ON kv.version_id = v.id\s*WHERE v.id = \$2 AND v.document_id = d.id AND v.source_id IS NOT NULL AND kv.decision = 'approved'\s*AND NOT EXISTS`).
		WithArgs(documentID, versionID).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), documentID, versionID); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition after newer validation, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishRejectsMissingAuthor(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	pool.ExpectExec(`UPDATE knowledge_documents d SET status = 'published', updated_at = NOW\(\) WHERE d.id = \$1 AND d.status = 'validated' AND NULLIF\(BTRIM\(d.author_name\), ''\) IS NOT NULL`).
		WithArgs("doc", "version").WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	if err := s.Publish(context.Background(), "doc", "version"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition without author, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAddVersionRejectsValidatedDocument(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1`).
		WithArgs("doc").
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("validated"))

	if err := s.AddVersion(context.Background(), "doc", 2, "content", nil); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition for validated document, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAddVersionRejectsPublishedDocument(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1`).
		WithArgs("doc").
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("published"))

	if err := s.AddVersion(context.Background(), "doc", 2, "content", nil); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition, got %v", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateApprovedOlderVersionKeepsDocumentInReview(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)
	tx, err := pool.ExpectBegin(), error(nil)
	_ = tx
	if err != nil {
		t.Fatal(err)
	}
	pool.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM knowledge_versions WHERE id = \$1 AND document_id = \$2\)`).
		WithArgs("version-1", "doc").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	pool.ExpectQuery(`SELECT status FROM knowledge_documents WHERE id = \$1`).
		WithArgs("doc").
		WillReturnRows(pgxmock.NewRows([]string{"status"}).AddRow("review"))
	pool.ExpectExec(`INSERT INTO knowledge_validations`).
		WithArgs("version-1", "validator", "approved", "").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	pool.ExpectQuery(`SELECT EXISTS\(\s*SELECT 1 FROM knowledge_versions v`).
		WithArgs("version-1", "doc").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
	pool.ExpectExec(`UPDATE knowledge_documents SET status = \$2, updated_at = NOW\(\) WHERE id = \$1`).
		WithArgs("doc", "review").WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	pool.ExpectCommit()
	if err := s.Validate(context.Background(), "doc", "version-1", "validator", "approved", ""); err != nil {
		t.Fatal(err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
