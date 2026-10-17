package knowledge

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *HTTPHandler) SearchPublished(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	items, err := h.service.SearchPublished(r.Context(), RetrievalFilter{
		CropID:        q.Get("crop_id"),
		GrowthStageID: q.Get("growth_stage_id"),
		Query:         q.Get("q"),
	}, limit)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}
