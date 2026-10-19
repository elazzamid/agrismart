package knowledge

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestRecommendFertilizersRequiresExpertOrAdmin(t *testing.T) {
	h := NewHTTPHandler(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge/fertilizers/recommendations?crop_id=c1&product_amount=10", nil)
	rec := httptest.NewRecorder()
	h.RecommendFertilizers(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestRecommendFertilizersRejectsInvalidProductAmount(t *testing.T) {
	cases := []string{"-1", "NaN", "+Inf", "-Inf"}
	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			h := NewHTTPHandler(nil)
			ctx := auth.WithRole(context.Background(), "expert")
			req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge/fertilizers/recommendations?crop_id=c1&product_amount="+value, nil).WithContext(ctx)
			rec := httptest.NewRecorder()
			h.RecommendFertilizers(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), "product_amount must be a finite non-negative number") {
				t.Fatalf("unexpected response: %s", rec.Body.String())
			}
		})
	}
}

func TestRecommendFertilizersRequiresCropID(t *testing.T) {
	h := NewHTTPHandler(nil)
	ctx := auth.WithRole(context.Background(), "admin")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/knowledge/fertilizers/recommendations?product_amount=10", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	h.RecommendFertilizers(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "crop_id is required") {
		t.Fatalf("unexpected response: %s", rec.Body.String())
	}
}
