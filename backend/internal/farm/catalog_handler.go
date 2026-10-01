package farm

import (
	"net/http"

	"github.com/elazzamid/agrismart/backend/internal/auth"
)

type CatalogHandler struct{ service *CatalogService }

func NewCatalogHandler(service *CatalogService) *CatalogHandler { return &CatalogHandler{service: service} }

func (h *CatalogHandler) ListCrops(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.UserIDFromContext(r.Context()); !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.service.ListCrops(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListVarieties(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.UserIDFromContext(r.Context()); !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.service.ListVarieties(r.Context(), r.PathValue("cropID"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CatalogHandler) ListGrowthStages(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.UserIDFromContext(r.Context()); !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.service.ListGrowthStages(r.Context(), r.PathValue("cropID"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, items)
}
