package knowledge

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type HTTPHandler struct{ service *Service }

func NewHTTPHandler(service *Service) *HTTPHandler { return &HTTPHandler{service: service} }

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func role(r *http.Request) string {
	value, _ := r.Context().Value(auth.RoleContextKey()).(string)
	return value
}

func requireRole(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	got := role(r)
	for _, want := range allowed {
		if got == want {
			return true
		}
	}
	jsonError(w, http.StatusForbidden, "forbidden")
	return false
}

func documentID(path, suffix string) string {
	id := strings.TrimPrefix(path, "/api/v1/knowledge/documents/")
	return strings.TrimSuffix(id, suffix)
}

func (h *HTTPHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	var in CreateDocumentInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
	d, err := h.service.CreateDocument(r.Context(), in.Slug, in.Title, in.Summary, in.AuthorName)
	if err != nil { jsonError(w, 400, err.Error()); return }
	jsonResponse(w, 201, d)
}

func (h *HTTPHandler) AddVersion(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	id := documentID(r.URL.Path, "/versions")
	var in struct {
		Version int `json:"version"`
		Content string `json:"content"`
		SourceID *string `json:"source_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
	if err := h.service.AddVersion(r.Context(), id, in.Version, in.Content, in.SourceID); err != nil {
		status := 400
		if errors.Is(err, ErrDocumentNotFound) { status = 404 }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, 201, map[string]string{"status": "version_created"})
}

func (h *HTTPHandler) Validate(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	id := documentID(r.URL.Path, "/validate")
	var in struct {
		VersionID string `json:"version_id"`
		ValidatorName string `json:"validator_name"`
		Decision string `json:"decision"`
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
	if err := h.service.Validate(r.Context(), id, in.VersionID, in.ValidatorName, in.Decision, in.Notes); err != nil {
		status := 400
		if errors.Is(err, ErrDocumentNotFound) { status = 404 }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, 200, map[string]string{"status": "validated"})
}

func (h *HTTPHandler) Publish(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "admin") { return }
	id := documentID(r.URL.Path, "/publish")
	var in struct { VersionID string `json:"version_id"` }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
	if strings.TrimSpace(in.VersionID) == "" { jsonError(w, 400, "version_id is required"); return }
	if err := h.service.Publish(r.Context(), id, in.VersionID); err != nil {
		status := 409
		if errors.Is(err, ErrDocumentNotFound) { status = 404 }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, 200, map[string]string{"status": "published"})
}
