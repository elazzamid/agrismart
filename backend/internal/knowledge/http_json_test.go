package knowledge

import (
	"encoding/json"
	"strings"
	"testing"

func TestDecodeSingleJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
		wantErr bool
	}{
		{name: "single object", body: "{\"value\":\"ok\"}"},
		{name: "single object with whitespace", body: "  {\"value\":\"ok\"} \n"},
		{name: "second JSON value", body: "{\"value\":\"ok\"} {\"value\":\"unexpected\"}", wantErr: true},
		{name: "trailing garbage", body: "{\"value\":\"ok\"} trailing", wantErr: true},
		{name: "empty body", body: "", wantErr: true},
		{name: "malformed JSON", body: "{\"value\":", wantErr: true},
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
