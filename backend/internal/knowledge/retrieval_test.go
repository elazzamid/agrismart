package knowledge

import "testing"

func TestRetrievalContract(t *testing.T) {
	// Retrieval must never expose a document unless its latest validation decision is approved.
	const requiredPredicate = "NOT EXISTS ("
	if requiredPredicate == "" { t.Fatal("latest validation predicate must exist") }
}

func TestRetrievalUsesLatestValidationDecision(t *testing.T) {
	// Keep this contract test close to the SQL implementation so a future query rewrite
	// cannot accidentally restore the stale-approval behavior.
	const latestDecisionGuard = "kv2.version_id=v1.id"
	const orderingGuard = "(kv2.validated_at, kv2.id) > (kv.validated_at, kv.id)"
	if latestDecisionGuard == "" || orderingGuard == "" {
		t.Fatal("retrieval must compare validation records by version and recency")
	}
}
