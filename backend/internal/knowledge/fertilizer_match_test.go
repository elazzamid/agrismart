package knowledge

import "testing"

func TestMatchFertilizerToNutrient(t *testing.T) {
    got, err := MatchFertilizerToNutrient("f1", "Urea", "N", 46, 10)
    if err != nil { t.Fatal(err) }
    if got.FertilizerID != "f1" || got.NutrientCode != "N" { t.Fatalf("unexpected match: %+v", got) }
    if got.CalculatedProductAmount < 21.7391 || got.CalculatedProductAmount > 21.7392 { t.Fatalf("unexpected amount: %v", got.CalculatedProductAmount) }
}
