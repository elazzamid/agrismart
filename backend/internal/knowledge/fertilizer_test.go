package knowledge

import "testing"

func TestCalculateFertilizerAmount(t *testing.T) {
    got, err := CalculateFertilizerAmount("f1", "Urea", "N", 10, 46, "kg")
    if err != nil { t.Fatal(err) }
    if got.RequiredFertilizerAmount < 21.7391 || got.RequiredFertilizerAmount > 21.7392 { t.Fatalf("unexpected amount: %v", got.RequiredFertilizerAmount) }
}

func TestCalculateFertilizerAmountRejectsInvalidPercentage(t *testing.T) {
    if _, err := CalculateFertilizerAmount("f1", "", "N", 10, 0, "kg"); err == nil { t.Fatal("expected error") }
    if _, err := CalculateFertilizerAmount("f1", "", "N", 10, 101, "kg"); err == nil { t.Fatal("expected error") }
}
