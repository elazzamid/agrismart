package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestDiagnosePublishedUsesLatestVersionOnly(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := NewService(pool)

	pool.ExpectQuery(`(?s)SELECT DISTINCT d\.id, d\.slug, d\.title,.*FROM knowledge_documents d\s*JOIN knowledge_versions v ON v\.document_id=d\.id.*WHERE d\.status='published'\s*AND v\.version_no = \(SELECT MAX\(v2\.version_no\) FROM knowledge_versions v2 WHERE v2\.document_id=d\.id\).*AND EXISTS \(\s*SELECT 1 FROM knowledge_validations kv\s*WHERE kv\.version_id=v\.id AND kv\.decision='approved'\s*AND NOT EXISTS \(\s*SELECT 1 FROM knowledge_validations kv2`).
		WithArgs("crop-1", []string{"leaf spot"}, 10).
		WillReturnRows(pgxmock.NewRows([]string{"id", "slug", "title", "problem_type", "problem_id", "score", "content"}).
			AddRow("doc-1", "cabai-leaf-spot", "Leaf Spot", "disease", "disease-1", 1, "leaf spot on leaves"))

	items, err := s.DiagnosePublished(context.Background(), "crop-1", []string{" leaf spot "}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].DocumentID != "doc-1" {
		t.Fatalf("unexpected diagnosis candidates: %+v", items)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
