package knowledge

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"

    "github.com/elazzamid/agrismart/backend/internal/auth"
)

type HTTPHandler struct { service *Service }
func NewHTTPHandler(service *Service) *HTTPHandler { return &HTTPHandler{service: service} }

func jsonResponse(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}
func jsonError(w http.ResponseWriter, status int, message string) { jsonResponse(w, status, map[string]string{"error": message}) }

func role(r *http.Request) string {
    value, _ := r.Context().Value(auth.RoleContextKey()).(string)
    return value
}
func requireRole(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
    got := role(r)
    for _, want := range allowed { if got == want { return true } }
    jsonError(w, http.StatusForbidden, "forbidden")
    return false
}

func (h *HTTPHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
    if !requireRole(w, r, "expert", "admin") { return }
    var in CreateDocumentInput
    dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
    d, err := h.service.CreateDocument(r.Context(), in)
    if err != nil { jsonError(w, 400, err.Error()); return }
    jsonResponse(w, 201, d)
}

func (h *HTTPHandler) AddVersion(w http.ResponseWriter, r *http.Request) {
    if !requireRole(w, r, "expert", "admin") { return }
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/knowledge/documents/")
    id = strings.TrimSuffix(id, "/versions")
    var in CreateVersionInput
    dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
    v, err := h.service.AddVersion(r.Context(), id, in)
    if err != nil { jsonError(w, 400, err.Error()); return }
    jsonResponse(w, 201, v)
}

func (h *HTTPHandler) Validate(w http.ResponseWriter, r *http.Request) {
    if !requireRole(w, r, "expert", "admin") { return }
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/knowledge/documents/")
    id = strings.TrimSuffix(id, "/validate")
    var in ValidationInput
    dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil { jsonError(w, 400, "invalid request body"); return }
    if err := h.service.Validate(r.Context(), id, in); err != nil {
        status := 400; if errors.Is(err, ErrNotFound) { status = 404 }
        jsonError(w, status, err.Error()); return
    }
    jsonResponse(w, 200, map[string]string{"status":"validated"})
}

func (h *HTTPHandler) Publish(w http.ResponseWriter, r *http.Request) {
    if !requireRole(w, r, "admin") { return }
    id := strings.TrimPrefix(r.URL.Path, "/api/v1/knowledge/documents/")
    id = strings.TrimSuffix(id, "/publish")
    if err := h.service.Publish(r.Context(), id); err != nil {
        status := 400; if errors.Is(err, ErrNotFound) { status = 404 }
        if errors.Is(err, ErrNotValidated) { status = 409 }
        jsonError(w, status, err.Error()); return
    }
    jsonResponse(w, 200, map[string]string{"status":"published"})
}
