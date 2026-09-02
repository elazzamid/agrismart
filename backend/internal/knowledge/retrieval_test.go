package knowledge

import (
	"strings"
	"testing"
)

func TestRetrievalContract(t *testing.T) {
	// Retrieval must never expose a document unless its latest validation decision is approved.
	for _, required := range []string{
		"d.status='published'",
		"v1.version_no = (SELECT MAX(v2.version_no) FROM knowledge_versions v2 WHERE v2.document_id=d.id)",
		"NULLIF(BTRIM(v1.source_id), '') IS NOT NULL",
		"kv.decision='approved'",
		"NOT EXISTS (",
	} {
		if !strings.Contains(publishedRetrievalQuery, required) {
			t.Fatalf("retrieval query missing required guard: %s", required)
		}
	}
}

func TestRetrievalUsesLatestValidationDecision(t *testing.T) {
	// Keep this contract test close to the SQL implementation so a future query rewrite
	// cannot accidentally restore stale-version or stale-approval behavior.
	for _, required := range []string{
		"kv.version_id=v1.id",
		"kv2.version_id=v1.id",
		"(kv2.validated_at, kv2.id) > (kv.validated_at, kv.id)",
	} {
		if !strings.Contains(publishedRetrievalQuery, required) {
			t.Fatalf("retrieval query missing latest-validation guard: %s", required)
		}
	}
}
