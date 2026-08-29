package knowledge

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestRecommendFertilizersRejectsNonFiniteProductAmount(t *testing.T) {
	h := NewHTTPHandler(nil)
	ctx := auth.WithRole(context.Background(), "expert")
	for _, amount := range []string{"NaN", "+Inf", "-Inf"} {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge/fertilizers/recommendations?crop_id=c1&product_amount="+amount, nil).WithContext(ctx)
		rec := httptest.NewRecorder()
		h.RecommendFertilizers(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %s, got %d", amount, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "product_amount must be a finite non-negative number") {
			t.Fatalf("unexpected response for %s: %s", amount, rec.Body.String())
		}
	}
}
