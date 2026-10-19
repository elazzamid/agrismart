package farm

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type CropCycleHandler struct{ service *CropCycleService }

func NewCropCycleHandler(service *CropCycleService) *CropCycleHandler {
	return &CropCycleHandler{service: service}
}

func (h *CropCycleHandler) List(w http.ResponseWriter, r *http.Request) {
	farmerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	cycles, err := h.service.List(r.Context(), farmerID, r.PathValue("farmID"), r.PathValue("plotID"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, cycles)
}

func (h *CropCycleHandler) Create(w http.ResponseWriter, r *http.Request) {
	farmerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var in CreateCropCycleInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := decodeSingleJSON(dec, &in); err != nil {
		status, message := http.StatusBadRequest, "invalid request body"
		if isRequestBodyTooLarge(err) {
			status, message = http.StatusRequestEntityTooLarge, "request body too large"
		}
		writeError(w, status, message)
		return
	}
	cycle, err := h.service.Create(r.Context(), farmerID, r.PathValue("farmID"), r.PathValue("plotID"), in)
	if errors.Is(err, ErrCropCycleNotFound) {
		writeError(w, http.StatusNotFound, "plot not found")
		return
	}
	if err != nil {
		if errors.Is(err, ErrInvalidCropCycleInput) { writeError(w, http.StatusBadRequest, err.Error()) } else { writeError(w, http.StatusInternalServerError, "internal server error") }
		return
	}
	writeJSON(w, http.StatusCreated, cycle)
}
