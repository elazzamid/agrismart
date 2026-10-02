package farm

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type Handler struct {
	service       *Service
	plotService   *PlotService
	catalogService *CatalogService
}

func NewHandler(service *Service, plotService *PlotService, catalogService *CatalogService) *Handler {
	return &Handler{service: service, plotService: plotService, catalogService: catalogService}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func userID(r *http.Request) (string, bool) { return auth.UserIDFromContext(r.Context()) }

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	farms, err := h.service.List(r.Context(), id)
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, farms)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	var in CreateInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { writeError(w, http.StatusBadRequest, "invalid request body"); return }
	f, err := h.service.Create(r.Context(), id, in)
	if err != nil { writeError(w, http.StatusBadRequest, err.Error()); return }
	writeJSON(w, http.StatusCreated, f)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	farmID := strings.TrimPrefix(r.URL.Path, "/api/v1/farms/")
	if farmID == "" { writeError(w, http.StatusBadRequest, "farm id is required"); return }
	f, err := h.service.Get(r.Context(), id, farmID)
	if errors.Is(err, ErrNotFound) { writeError(w, http.StatusNotFound, "farm not found"); return }
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, f)
}

func (h *Handler) ListPlots(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	farmID := strings.TrimPrefix(r.URL.Path, "/api/v1/farms/")
	farmID = strings.TrimSuffix(farmID, "/plots")
	farmID = strings.TrimSuffix(farmID, "/")
	if farmID == "" { writeError(w, http.StatusBadRequest, "farm id is required"); return }
	plots, err := h.plotService.List(r.Context(), id, farmID)
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, plots)
}

func (h *Handler) CreatePlot(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	farmID := strings.TrimPrefix(r.URL.Path, "/api/v1/farms/")
	farmID = strings.TrimSuffix(farmID, "/plots")
	farmID = strings.TrimSuffix(farmID, "/")
	if farmID == "" { writeError(w, http.StatusBadRequest, "farm id is required"); return }
	var in CreatePlotInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil { writeError(w, http.StatusBadRequest, "invalid request body"); return }
	plot, err := h.plotService.Create(r.Context(), id, farmID, in)
	if errors.Is(err, ErrPlotNotFound) { writeError(w, http.StatusNotFound, "farm not found"); return }
	if err != nil { writeError(w, http.StatusBadRequest, err.Error()); return }
	writeJSON(w, http.StatusCreated, plot)
}

func (h *Handler) ListCrops(w http.ResponseWriter, r *http.Request) {
	if _, ok := userID(r); !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	crops, err := h.catalogService.ListCrops(r.Context())
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, crops)
}

func (h *Handler) GetCrop(w http.ResponseWriter, r *http.Request) {
	if _, ok := userID(r); !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	cropID := strings.TrimPrefix(r.URL.Path, "/api/v1/crops/")
	if cropID == "" { writeError(w, http.StatusBadRequest, "crop id is required"); return }
	crop, err := h.catalogService.GetCrop(r.Context(), cropID)
	if errors.Is(err, ErrNotFound) { writeError(w, http.StatusNotFound, "crop not found"); return }
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, crop)
}

func (h *Handler) ListVarieties(w http.ResponseWriter, r *http.Request) {
	if _, ok := userID(r); !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	cropID := strings.TrimPrefix(r.URL.Path, "/api/v1/crops/")
	cropID = strings.TrimSuffix(cropID, "/varieties")
	cropID = strings.TrimSuffix(cropID, "/")
	if cropID == "" { writeError(w, http.StatusBadRequest, "crop id is required"); return }
	varieties, err := h.catalogService.ListVarieties(r.Context(), cropID)
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, varieties)
}

func (h *Handler) ListGrowthStages(w http.ResponseWriter, r *http.Request) {
	if _, ok := userID(r); !ok { writeError(w, http.StatusUnauthorized, "unauthorized"); return }
	cropID := strings.TrimPrefix(r.URL.Path, "/api/v1/crops/")
	cropID = strings.TrimSuffix(cropID, "/growth-stages")
	cropID = strings.TrimSuffix(cropID, "/")
	if cropID == "" { writeError(w, http.StatusBadRequest, "crop id is required"); return }
	stages, err := h.catalogService.ListGrowthStages(r.Context(), cropID)
	if err != nil { writeError(w, http.StatusInternalServerError, "internal server error"); return }
	writeJSON(w, http.StatusOK, stages)
}
