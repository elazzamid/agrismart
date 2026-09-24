package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct { service *Service }
func NewHandler(service *Service) *Handler { return &Handler{service: service} }

type credentialsRequest struct { Email string `json:"email"`; Password string `json:"password"` }
type registerRequest struct { Email string `json:"email"`; Password string `json:"password"`; FullName string `json:"full_name"`; Phone string `json:"phone"` }
type userResponse struct { ID string `json:"id"`; Email string `json:"email"`; Role string `json:"role"` }

func writeJSON(w http.ResponseWriter, status int, value any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(value) }
func writeError(w http.ResponseWriter, status int, message string) { writeJSON(w, status, map[string]string{"error": message}) }
func decodeJSON(r *http.Request, dst any) error { decoder := json.NewDecoder(r.Body); decoder.DisallowUnknownFields(); return decoder.Decode(dst) }

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil { writeError(w, http.StatusBadRequest, "invalid request body"); return }
	user, err := h.service.Register(r.Context(), req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		switch { case errors.Is(err, ErrEmailExists): writeError(w, http.StatusConflict, err.Error()); case errors.Is(err, ErrInvalidEmail), errors.Is(err, ErrWeakPassword): writeError(w, http.StatusBadRequest, err.Error()); default: writeError(w, http.StatusInternalServerError, "internal server error") }
		return
	}
	writeJSON(w, http.StatusCreated, userResponse{ID: user.ID, Email: user.Email, Role: user.Role})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := decodeJSON(r, &req); err != nil { writeError(w, http.StatusBadRequest, "invalid request body"); return }
	user, token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil { if errors.Is(err, ErrInvalidCredentials) { writeError(w, http.StatusUnauthorized, err.Error()); return }; writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, map[string]any{"access_token": token, "token_type": "Bearer", "user": userResponse{ID: user.ID, Email: user.Email, Role: user.Role}})
}

func (h *Handler) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") { writeError(w, http.StatusUnauthorized, "missing bearer token"); return }
		userID, _, err := h.service.tokens.Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil { writeError(w, http.StatusUnauthorized, "invalid token"); return }
		next.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), userID)))
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	user, err := h.service.UserByID(r.Context(), userID)
	if err != nil { writeError(w, http.StatusUnauthorized, "user not found"); return }
	writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Email: user.Email, Role: user.Role})
}
