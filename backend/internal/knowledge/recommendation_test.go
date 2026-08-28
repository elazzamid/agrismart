package knowledge

import "testing"

func TestBuildFertilizerRecommendationRequiresProvenance(t *testing.T) {
    calc, err := CalculateFertilizerAmount("f1", "Urea", "N", 10, 46, "kg")
    if err != nil { t.Fatal(err) }
    if _, err := BuildFertilizerRecommendation(calc, "", "v1", false); err == nil { t.Fatal("expected missing source error") }
}

func TestBuildFertilizerRecommendationDoesNotInferAgronomicRate(t *testing.T) {
    calc, err := CalculateFertilizerAmount("f1", "Urea", "N", 10, 46, "kg")
    if err != nil { t.Fatal(err) }
    rec, err := BuildFertilizerRecommendation(calc, "doc1", "v1", false)
    if err != nil { t.Fatal(err) }
    if rec.AgronomicRateSupported { t.Fatal("expected unsupported agronomic rate") }
}
