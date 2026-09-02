package knowledge

import (
    "encoding/json"
    "net/http"
)

type DiagnosisRequest struct {
    CropID string `json:"crop_id"`
    SymptomIDs []string `json:"symptom_ids"`
    Limit int `json:"limit"`
}

func (h *HTTPHandler) Diagnose(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { jsonError(w, http.StatusMethodNotAllowed, "method not allowed"); return }
    var in DiagnosisRequest
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil { jsonError(w, http.StatusBadRequest, "invalid request body"); return }
    result, err := h.service.DiagnoseBySymptoms(r.Context(), in.SymptomIDs, in.CropID, in.Limit)
    if err != nil { jsonError(w, http.StatusInternalServerError, "internal server error"); return }
    jsonResponse(w, http.StatusOK, map[string]any{"candidates": result, "disclaimer": "Hasil ini adalah kandidat berbasis gejala, bukan diagnosis pasti."})
}
