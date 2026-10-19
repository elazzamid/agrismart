package farm

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

var activityUUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type ActivityHandler struct {
	service *FarmActivityService
}

func NewActivityHandler(service *FarmActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) cropCycleID(r *http.Request) (string, bool) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/crop-cycles/")
	id = strings.TrimSuffix(id, "/activities")
	if id == "" || id == path || !activityUUIDPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cycleID, ok := h.cropCycleID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "valid crop cycle id is required")
		return
	}

	activities, err := h.service.List(r.Context(), userID, cycleID)
	if errors.Is(err, ErrFarmActivityNotFound) {
		writeError(w, http.StatusNotFound, "crop cycle not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, activities)
}

func (h *ActivityHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cycleID, ok := h.cropCycleID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "valid crop cycle id is required")
		return
	}

	var in CreateFarmActivityInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	activity, err := h.service.Create(r.Context(), userID, cycleID, in)
	switch {
	case errors.Is(err, ErrFarmActivityNotFound):
		writeError(w, http.StatusNotFound, "crop cycle not found")
	case errors.Is(err, ErrInvalidFarmActivityInput):
		writeError(w, http.StatusBadRequest, "activity_type, valid activity_date, and notes up to 2000 characters are required")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal server error")
	default:
		writeJSON(w, http.StatusCreated, activity)
	}
}
