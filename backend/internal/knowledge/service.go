package knowledge

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDocumentNotFound  = errors.New("knowledge document not found")
	ErrInvalidTransition = errors.New("invalid knowledge publication transition")
)

type Document struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Summary    string `json:"summary,omitempty"`
	Status     string `json:"status"`
	AuthorName string `json:"author_name,omitempty"`
}

type CreateDocumentInput struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Summary    string `json:"summary,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
}

type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct{ db DBTX }

func NewService(db *pgxpool.Pool) *Service { return &Service{db: db} }

func (s *Service) CreateDocument(ctx context.Context, slug, title, summary, author string) (Document, error) {
	slug, title = strings.TrimSpace(slug), strings.TrimSpace(title)
	if slug == "" || title == "" {
		return Document{}, errors.New("slug and title are required")
	}
	var d Document
	err := s.db.QueryRow(ctx, `
		INSERT INTO knowledge_documents (slug, title, summary, author_name)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''))
		RETURNING id, slug, title, COALESCE(summary, ''), status, COALESCE(author_name, '')`, slug, title, strings.TrimSpace(summary), strings.TrimSpace(author)).
		Scan(&d.ID, &d.Slug, &d.Title, &d.Summary, &d.Status, &d.AuthorName)
	return d, err
}

func (s *Service) AddVersion(ctx context.Context, documentID string, version int, content string, sourceID *string) error {
	if version < 1 || strings.TrimSpace(content) == "" {
		return errors.New("version and content are required")
	}
	var status string
	if err := s.db.QueryRow(ctx, `SELECT status FROM knowledge_documents WHERE id = $1`, documentID).Scan(&status); errors.Is(err, pgx.ErrNoRows) {
		return ErrDocumentNotFound
	} else if err != nil {
		return err
	}
	if status == "published" || status == "archived" {
		return ErrInvalidTransition
	}
	_, err := s.db.Exec(ctx, `INSERT INTO knowledge_versions (document_id, version_no, content, source_id) VALUES ($1, $2, $3, $4)`, documentID, version, strings.TrimSpace(content), sourceID)
	return err
}

func (s *Service) Validate(ctx context.Context, documentID string, versionID, validator, decision, notes string) error {
	validator = strings.TrimSpace(validator)
	if validator == "" {
		return errors.New("validator is required")
	}
	if decision != "approved" && decision != "rejected" && decision != "needs_revision" {
		return errors.New("invalid validation decision")
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM knowledge_versions WHERE id = $1 AND document_id = $2)`, versionID, documentID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrDocumentNotFound
	}
	_, err = tx.Exec(ctx, `INSERT INTO knowledge_validations (version_id, validator_name, decision, notes) VALUES ($1, $2, $3, NULLIF($4, ''))`, versionID, validator, decision, strings.TrimSpace(notes))
	if err != nil {
		return err
	}
	status := "review"
	if decision == "approved" {
		status = "validated"
	}
	if decision == "rejected" {
		status = "draft"
	}
	if _, err = tx.Exec(ctx, `UPDATE knowledge_documents SET status = $2, updated_at = NOW() WHERE id = $1`, documentID, status); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) Publish(ctx context.Context, documentID, versionID string) error {
	result, err := s.db.Exec(ctx, `UPDATE knowledge_documents d SET status = 'published', updated_at = NOW() WHERE d.id = $1 AND d.status = 'validated' AND EXISTS (
		SELECT 1 FROM knowledge_versions v
		JOIN knowledge_validations kv ON kv.version_id = v.id
		WHERE v.id = $2 AND v.document_id = d.id AND kv.decision = 'approved'
		AND v.version_no = (SELECT MAX(v2.version_no) FROM knowledge_versions v2 WHERE v2.document_id = d.id)
	)`, documentID, versionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrInvalidTransition
	}
	return nil
}
