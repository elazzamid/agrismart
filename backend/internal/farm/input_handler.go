package farm

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

var farmInputUUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type FarmInputHandler struct {
	service *FarmInputService
}

func NewFarmInputHandler(service *FarmInputService) *FarmInputHandler {
	return &FarmInputHandler{service: service}
}

func (h *FarmInputHandler) cropCycleID(r *http.Request) (string, bool) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/crop-cycles/")
	id = strings.TrimSuffix(id, "/inputs")
	if id == "" || id == path || !farmInputUUIDPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

func (h *FarmInputHandler) List(w http.ResponseWriter, r *http.Request) {
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

	inputs, err := h.service.List(r.Context(), userID, cycleID)
	if errors.Is(err, ErrFarmInputNotFound) {
		writeError(w, http.StatusNotFound, "crop cycle not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, inputs)
}

func (h *FarmInputHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var in CreateFarmInputInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input, err := h.service.Create(r.Context(), userID, cycleID, in)
	switch {
	case errors.Is(err, ErrFarmInputNotFound):
		writeError(w, http.StatusNotFound, "crop cycle not found")
	case errors.Is(err, ErrInvalidFarmInput):
		writeError(w, http.StatusBadRequest, "input_type, product_name, positive quantity, unit, and valid applied_date are required")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal server error")
	default:
		writeJSON(w, http.StatusCreated, input)
	}
}
