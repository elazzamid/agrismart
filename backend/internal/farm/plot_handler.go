package farm

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type PlotHandler struct{ service *PlotService }

func NewPlotHandler(service *PlotService) *PlotHandler { return &PlotHandler{service: service} }

func (h *PlotHandler) List(w http.ResponseWriter, r *http.Request) {
	farmerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	plots, err := h.service.List(r.Context(), farmerID, r.PathValue("farmID"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, plots)
}

func (h *PlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	farmerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var in CreatePlotInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	plot, err := h.service.Create(r.Context(), farmerID, r.PathValue("farmID"), in)
	if errors.Is(err, ErrPlotNotFound) {
		writeError(w, http.StatusNotFound, "farm not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, plot)
}
