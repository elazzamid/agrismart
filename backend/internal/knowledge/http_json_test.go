package knowledge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestDecodeSingleJSON(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "single object", body: `{"value":"ok"}`},
		{name: "single object with whitespace", body: "  {\"value\":\"ok\"} \n"},
		{name: "second JSON value", body: `{"value":"ok"} {"value":"unexpected"}`, wantErr: true},
		{name: "trailing garbage", body: `{"value":"ok"} trailing`, wantErr: true},
		{name: "empty body", body: "", wantErr: true},
		{name: "malformed JSON", body: `{"value":`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got struct {
				Value string `json:"value"`
			}
			err := decodeSingleJSON(json.NewDecoder(strings.NewReader(tt.body)), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeSingleJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.Value != "ok" {
				t.Fatalf("decoded value = %q, want %q", got.Value, "ok")
			}
		})
	}
}

func TestKnowledgeHandlersRejectTrailingJSON(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(*HTTPHandler, http.ResponseWriter, *http.Request)
		role string
	}{
		{
			name: "create document",
			path: "/api/v1/knowledge/documents",
			body: `{"slug":"test","title":"Test","summary":"Summary","author_name":"Expert"} {}`,
			call: (*HTTPHandler).CreateDocument,
			role: "expert",
		},
		{
			name: "diagnosis",
			path: "/api/v1/diagnosis",
			body: `{"crop_id":"crop-1","symptom_ids":["symptom-1"]} {}`,
			call: (*HTTPHandler).Diagnose,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHTTPHandler(nil)
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			if tt.role != "" {
				req = req.WithContext(auth.WithRole(req.Context(), tt.role))
			}
			rec := httptest.NewRecorder()

			tt.call(h, rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d for trailing JSON, got %d; body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
			}
		})
	}
}
