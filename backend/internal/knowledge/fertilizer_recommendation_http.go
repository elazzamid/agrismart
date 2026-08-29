package knowledge

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// RecommendFertilizers exposes source-backed fertilizer candidates and nutrient
// coverage. It does not infer or return agronomic application rates.
func (h *HTTPHandler) RecommendFertilizers(w http.ResponseWriter, r *http.Request) {
	if !requireRole(w, r, "expert", "admin") { return }
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	productAmount, err := strconv.ParseFloat(q.Get("product_amount"), 64)
	if err != nil || productAmount < 0 {
		jsonError(w, http.StatusBadRequest, "product_amount must be a non-negative number")
		return
	}
	cropID := q.Get("crop_id")
	growthStageID := q.Get("growth_stage_id")
	if cropID == "" {
		jsonError(w, http.StatusBadRequest, "crop_id is required")
		return
	}
	items, err := h.service.BuildFertilizerCandidates(r.Context(), cropID, growthStageID, productAmount, limit)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}
