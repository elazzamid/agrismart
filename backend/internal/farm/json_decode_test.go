package farm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

func TestDecodeSingleJSONRejectsTrailingData(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{name: "single JSON value", body: `{"name":"field"}`},
		{name: "whitespace after JSON", body: "  {\"name\":\"field\"} \n"},
		{name: "second JSON value", body: `{"name":"field"} {}`, wantErr: true},
		{name: "trailing garbage", body: `{"name":"field"} trailing`, wantErr: true},
		{name: "empty body", body: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]string
			err := decodeSingleJSON(json.NewDecoder(strings.NewReader(tt.body)), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeSingleJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFarmMutationHandlersRejectTrailingJSON(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(http.ResponseWriter, *http.Request)
	}{
		{
			name: "create farm",
			path: "/api/v1/farms",
			body: `{"name":"Field"} {}`,
			call: func(w http.ResponseWriter, r *http.Request) { NewHandler(nil).Create(w, r) },
		},
		{
			name: "create plot",
			path: "/api/v1/farms/farm-1/plots",
			body: `{"name":"Plot"} {}`,
			call: func(w http.ResponseWriter, r *http.Request) { NewPlotHandler(nil).Create(w, r) },
		},
		{
			name: "create crop cycle",
			path: "/api/v1/farms/farm-1/plots/plot-1/crop-cycles",
			body: `{"crop_id":"crop-1"} {}`,
			call: func(w http.ResponseWriter, r *http.Request) { NewCropCycleHandler(nil).Create(w, r) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req = req.WithContext(auth.WithUserID(req.Context(), "farmer-1"))
			rec := httptest.NewRecorder()

			tt.call(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d for trailing JSON, got %d; body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
			}
		})
	}
}
