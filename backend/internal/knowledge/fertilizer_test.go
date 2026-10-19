package knowledge

import (
    "math"
    "testing"
)

func TestCalculateFertilizerAmount(t *testing.T) {
    got, err := CalculateFertilizerAmount("f1", "Urea", "N", 10, 46, "kg")
    if err != nil { t.Fatal(err) }
    if got.RequiredFertilizerAmount < 21.7391 || got.RequiredFertilizerAmount > 21.7392 { t.Fatalf("unexpected amount: %v", got.RequiredFertilizerAmount) }
}

func TestCalculateFertilizerAmountRejectsInvalidPercentage(t *testing.T) {
    if _, err := CalculateFertilizerAmount("f1", "", "N", 10, 0, "kg"); err == nil { t.Fatal("expected error") }
    if _, err := CalculateFertilizerAmount("f1", "", "N", 10, 101, "kg"); err == nil { t.Fatal("expected error") }
}

func TestCalculateFertilizerAmountRejectsNonFiniteInputs(t *testing.T) {
    cases := []struct {
        name       string
        target     float64
        percentage float64
    }{
        {name: "nan target", target: math.NaN(), percentage: 46},
        {name: "infinite target", target: math.Inf(1), percentage: 46},
        {name: "nan percentage", target: 10, percentage: math.NaN()},
        {name: "infinite percentage", target: 10, percentage: math.Inf(1)},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            if _, err := CalculateFertilizerAmount("f1", "Urea", "N", tc.target, tc.percentage, "kg"); err == nil {
                t.Fatal("expected error for non-finite input")
            }
        })
    }
}

func TestCalculateFertilizerAmountRejectsOverflow(t *testing.T) {
    if _, err := CalculateFertilizerAmount("f1", "Urea", "N", math.MaxFloat64, math.SmallestNonzeroFloat64, "kg"); err == nil {
        t.Fatal("expected error for non-finite calculated amount")
    }
}
