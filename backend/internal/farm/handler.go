package farm

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type Handler struct { service *Service }
func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) { writeJSON(w, status, map[string]string{"error": message}) }

func userID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(auth.UserContextKey()).(string)
	return id, ok && id != ""
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r); if !ok { writeError(w, 401, "unauthorized"); return }
	farms, err := h.service.List(r.Context(), id); if err != nil { writeError(w, 500, "internal server error"); return }
	writeJSON(w, 200, farms)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r); if !ok { writeError(w, 401, "unauthorized"); return }
	var in CreateInput
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { writeError(w, 400, "invalid request body"); return }
	f, err := h.service.Create(r.Context(), id, in)
	if err != nil { writeError(w, 400, err.Error()); return }
	writeJSON(w, 201, f)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r); if !ok { writeError(w, 401, "unauthorized"); return }
	farmID := strings.TrimPrefix(r.URL.Path, "/api/v1/farms/")
	if farmID == "" { writeError(w, 400, "farm id is required"); return }
	f, err := h.service.Get(r.Context(), id, farmID)
	if err == ErrNotFound { writeError(w, 404, "farm not found"); return }
	if err != nil { writeError(w, 500, "internal server error"); return }
	writeJSON(w, 200, f)
}
