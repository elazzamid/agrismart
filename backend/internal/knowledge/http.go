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

func jsonResponse(w http.ResponseWriter, status int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v) }
func jsonError(w http.ResponseWriter, status int, message string) { jsonResponse(w, status, map[string]string{"error": message}) }

func requireRole(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	got, ok := auth.RoleFromContext(r.Context())
	if !ok { jsonError(w, http.StatusForbidden, "forbidden"); return false }
	for _, want := range allowed { if got == want { return true } }
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
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { jsonError(w, http.StatusBadRequest, "invalid request body"); return }
	d, err := h.service.CreateDocument(r.Context(), in.Slug, in.Title, in.Summary, in.AuthorName)
	if err != nil { jsonError(w, http.StatusBadRequest, err.Error()); return }
	jsonResponse(w, http.StatusCreated, d)
}

func (h *HTTPHandler) AddVersion(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	id := documentID(r.URL.Path, "/versions")
	var in struct { Version int `json:"version"`; Content string `json:"content"`; SourceID *string `json:"source_id"` }
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { jsonError(w, http.StatusBadRequest, "invalid request body"); return }
	if err := h.service.AddVersion(r.Context(), id, in.Version, in.Content, in.SourceID); err != nil {
		status := http.StatusBadRequest; if errors.Is(err, ErrDocumentNotFound) { status = http.StatusNotFound }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, http.StatusCreated, map[string]string{"status": "version_created"})
}

func (h *HTTPHandler) Validate(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	id := documentID(r.URL.Path, "/validate")
	var in struct { VersionID string `json:"version_id"`; ValidatorName string `json:"validator_name"`; Decision string `json:"decision"`; Notes string `json:"notes"` }
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { jsonError(w, http.StatusBadRequest, "invalid request body"); return }
	if err := h.service.Validate(r.Context(), id, in.VersionID, in.ValidatorName, in.Decision, in.Notes); err != nil {
		status := http.StatusBadRequest; if errors.Is(err, ErrDocumentNotFound) { status = http.StatusNotFound }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "validated"})
}

func (h *HTTPHandler) Publish(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "admin") { return }
	id := documentID(r.URL.Path, "/publish")
	var in struct { VersionID string `json:"version_id"` }
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { jsonError(w, http.StatusBadRequest, "invalid request body"); return }
	if strings.TrimSpace(in.VersionID) == "" { jsonError(w, http.StatusBadRequest, "version_id is required"); return }
	if err := h.service.Publish(r.Context(), id, in.VersionID); err != nil {
		status := http.StatusConflict; if errors.Is(err, ErrDocumentNotFound) { status = http.StatusNotFound }
		jsonError(w, status, err.Error()); return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"status": "published"})
}
