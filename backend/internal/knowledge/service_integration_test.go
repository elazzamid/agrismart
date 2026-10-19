package knowledge

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// These tests exercise PostgreSQL row locks and transaction behavior. They are
// skipped for local unit-test runs unless DATABASE_URL points to a migrated DB.
func integrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set; skipping PostgreSQL integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open PostgreSQL pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("ping PostgreSQL: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func seedLifecycleDocument(t *testing.T, pool *pgxpool.Pool, status string, versionCount int, approveFirst bool) (string, []string) {
	t.Helper()
	ctx := context.Background()
	slug := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	var sourceID, documentID string
	if err := pool.QueryRow(ctx, `INSERT INTO knowledge_sources (title, source_type) VALUES ($1, 'reference') RETURNING id::text`, slug+" source").Scan(&sourceID); err != nil {
		t.Fatalf("insert source: %v", err)
	}
	if err := pool.QueryRow(ctx, `INSERT INTO knowledge_documents (slug, title, author_name, status) VALUES ($1, $2, 'integration test', $3) RETURNING id::text`, slug, slug, status).Scan(&documentID); err != nil {
		t.Fatalf("insert document: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM knowledge_documents WHERE id = $1`, documentID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM knowledge_sources WHERE id = $1`, sourceID)
	})
	versionIDs := make([]string, 0, versionCount)
	for i := 1; i <= versionCount; i++ {
		var versionID string
		if err := pool.QueryRow(ctx, `INSERT INTO knowledge_versions (document_id, version_no, content, source_id) VALUES ($1, $2, $3, $4) RETURNING id::text`, documentID, i, fmt.Sprintf("content version %d", i), sourceID).Scan(&versionID); err != nil {
			t.Fatalf("insert version %d: %v", i, err)
		}
		versionIDs = append(versionIDs, versionID)
		if i == 1 && approveFirst {
			if _, err := pool.Exec(ctx, `INSERT INTO knowledge_validations (version_id, validator_name, decision) VALUES ($1, 'integration reviewer', 'approved')`, versionID); err != nil {
				t.Fatalf("insert approved validation: %v", err)
			}
		}
	}
	return documentID, versionIDs
}

func TestPostgresValidateOlderVersionPreservesLatestStatus(t *testing.T) {
	pool := integrationPool(t)
	documentID, versions := seedLifecycleDocument(t, pool, "validated", 2, true)

	service := NewService(pool)
	if err := service.Validate(context.Background(), documentID, versions[0], "second reviewer", "rejected", "historical review"); err != nil {
		t.Fatalf("Validate(older version): %v", err)
	}

	var status string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM knowledge_documents WHERE id = $1`, documentID).Scan(&status); err != nil {
		t.Fatalf("read document status: %v", err)
	}
	if status != "validated" {
		t.Fatalf("older-version validation changed document status: got %q, want %q", status, "validated")
	}
}

func TestPostgresPublishRejectsStaleVersion(t *testing.T) {
	pool := integrationPool(t)
	documentID, versions := seedLifecycleDocument(t, pool, "validated", 2, true)

	err := NewService(pool).Publish(context.Background(), documentID, versions[0])
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Publish(stale version) error = %v, want ErrInvalidTransition", err)
	}
	var status string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM knowledge_documents WHERE id = $1`, documentID).Scan(&status); err != nil {
		t.Fatalf("read document status: %v", err)
	}
	if status != "validated" {
		t.Fatalf("rejected stale publish changed status: got %q, want %q", status, "validated")
	}
}

func TestPostgresConcurrentAddVersionAndPublishSerialize(t *testing.T) {
	pool := integrationPool(t)
	documentID, versions := seedLifecycleDocument(t, pool, "validated", 1, true)
	service := NewService(pool)

	start := make(chan struct{})
	var wg sync.WaitGroup
	var addErr, publishErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		addErr = service.AddVersion(context.Background(), documentID, 2, "concurrent version", nil)
	}()
	go func() {
		defer wg.Done()
		<-start
		publishErr = service.Publish(context.Background(), documentID, versions[0])
	}()
	close(start)
	wg.Wait()

	if publishErr != nil {
		t.Fatalf("Publish() error = %v, want success", publishErr)
	}
	if !errors.Is(addErr, ErrInvalidTransition) {
		t.Fatalf("AddVersion() error = %v, want ErrInvalidTransition", addErr)
	}
	var status string
	var versionCount int
	if err := pool.QueryRow(context.Background(), `SELECT d.status, (SELECT count(*) FROM knowledge_versions v WHERE v.document_id = d.id) FROM knowledge_documents d WHERE d.id = $1`, documentID).Scan(&status, &versionCount); err != nil {
		t.Fatalf("read final lifecycle state: %v", err)
	}
	if status != "published" || versionCount != 1 {
		t.Fatalf("unexpected final state: status=%q versions=%d; want published and 1 version", status, versionCount)
	}
}
