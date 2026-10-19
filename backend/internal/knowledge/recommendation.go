package knowledge

import "errors"

type FertilizerRecommendation struct {
    FertilizerID string `json:"fertilizer_id"`
    FertilizerName string `json:"fertilizer_name"`
    NutrientCode string `json:"nutrient_code"`
    TargetAmount float64 `json:"target_amount"`
    Unit string `json:"unit"`
    CalculatedProductAmount float64 `json:"calculated_product_amount"`
    SourceDocumentID string `json:"source_document_id"`
    SourceVersionID string `json:"source_version_id"`
    AgronomicRateSupported bool `json:"agronomic_rate_supported"`
    Disclaimer string `json:"disclaimer"`
}

// BuildFertilizerRecommendation turns a source-backed nutrient requirement and
// composition calculation into a recommendation-shaped response. An agronomic
// application rate is never inferred from nutrient arithmetic alone.
func BuildFertilizerRecommendation(calculation FertilizerAmount, sourceDocumentID, sourceVersionID string, agronomicRateSupported bool) (FertilizerRecommendation, error) {
    if sourceDocumentID == "" || sourceVersionID == "" { return FertilizerRecommendation{}, errors.New("source document and version are required") }
    if calculation.RequiredFertilizerAmount < 0 { return FertilizerRecommendation{}, errors.New("calculated amount cannot be negative") }
    disclaimer := "Jumlah produk adalah hasil konversi kandungan nutrisi; bukan otomatis dosis aplikasi agronomi."
    if agronomicRateSupported { disclaimer = "Dosis aplikasi hanya boleh digunakan sesuai sumber agronomi/label yang tervalidasi dan kondisi lapangan." }
    return FertilizerRecommendation{FertilizerID: calculation.FertilizerID, FertilizerName: calculation.FertilizerName, NutrientCode: calculation.NutrientCode, TargetAmount: calculation.TargetAmount, Unit: calculation.Unit, CalculatedProductAmount: calculation.RequiredFertilizerAmount, SourceDocumentID: sourceDocumentID, SourceVersionID: sourceVersionID, AgronomicRateSupported: agronomicRateSupported, Disclaimer: disclaimer}, nil
}
