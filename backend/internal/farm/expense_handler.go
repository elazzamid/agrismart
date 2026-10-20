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

var farmExpenseUUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type FarmExpenseHandler struct {
	service *FarmExpenseService
}

func NewFarmExpenseHandler(service *FarmExpenseService) *FarmExpenseHandler {
	return &FarmExpenseHandler{service: service}
}

func (h *FarmExpenseHandler) cropCycleID(r *http.Request) (string, bool) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/crop-cycles/")
	id = strings.TrimSuffix(id, "/expenses")
	if id == "" || id == path || !farmExpenseUUIDPattern.MatchString(id) {
		return "", false
	}
	return id, true
}

func (h *FarmExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
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

	expenses, err := h.service.List(r.Context(), userID, cycleID)
	if errors.Is(err, ErrFarmExpenseNotFound) {
		writeError(w, http.StatusNotFound, "crop cycle not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *FarmExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var in CreateFarmExpenseInput
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

	expense, err := h.service.Create(r.Context(), userID, cycleID, in)
	switch {
	case errors.Is(err, ErrFarmExpenseNotFound):
		writeError(w, http.StatusNotFound, "crop cycle not found")
	case errors.Is(err, ErrInvalidFarmExpenseInput):
		writeError(w, http.StatusBadRequest, "expense_type, description, positive amount, and valid expense_date are required")
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal server error")
	default:
		writeJSON(w, http.StatusCreated, expense)
	}
}
