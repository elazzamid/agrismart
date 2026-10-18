package knowledge

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestKnowledgeMutationEndpointsRejectFarmerRole(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		call   func(*HTTPHandler, http.ResponseWriter, *http.Request)
	}{
		{
			name: "create document", method: http.MethodPost, path: "/api/v1/knowledge/documents",
			call: (*HTTPHandler).CreateDocument,
		},
		{
			name: "add version", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/versions",
			call: (*HTTPHandler).AddVersion,
		},
		{
			name: "validate document", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/validate",
			call: (*HTTPHandler).Validate,
		},
		{
			name: "publish document", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/publish",
			call: (*HTTPHandler).Publish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(nil)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(auth.WithRole(req.Context(), "farmer"))
			rec := httptest.NewRecorder()

			tt.call(h, rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected status %d for farmer role, got %d; body=%s", http.StatusForbidden, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestKnowledgeMutationEndpointsRejectMissingRole(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		call   func(*HTTPHandler, http.ResponseWriter, *http.Request)
	}{
		{
			name: "create document", method: http.MethodPost, path: "/api/v1/knowledge/documents",
			call: (*HTTPHandler).CreateDocument,
		},
		{
			name: "add version", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/versions",
			call: (*HTTPHandler).AddVersion,
		},
		{
			name: "validate document", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/validate",
			call: (*HTTPHandler).Validate,
		},
		{
			name: "publish document", method: http.MethodPost, path: "/api/v1/knowledge/documents/doc-1/publish",
			call: (*HTTPHandler).Publish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(nil)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			tt.call(h, rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected status %d without role, got %d; body=%s", http.StatusForbidden, rec.Code, rec.Body.String())
			}
		})
	}
}
