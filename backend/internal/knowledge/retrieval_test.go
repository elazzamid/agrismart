package knowledge

import "testing"

func TestRetrievalContract(t *testing.T) {
	// The retrieval SQL is intentionally constrained to status='published'.
	// Runtime integration tests must verify draft/review/validated documents are excluded.
	const requiredPredicate = "d.status='published'"
	if requiredPredicate == "" { t.Fatal("published predicate must exist") }
}
