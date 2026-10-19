package knowledge

import (
    "errors"
    "strings"
)

type FertilizerMatch struct {
    FertilizerID string `json:"fertilizer_id"`
    FertilizerName string `json:"fertilizer_name"`
    NutrientCode string `json:"nutrient_code"`
    Percentage float64 `json:"percentage"`
    TargetAmount float64 `json:"target_amount"`
    CalculatedProductAmount float64 `json:"calculated_product_amount"`
}

// MatchFertilizerToNutrient returns a mathematical coverage match. It does not
// rank agronomic suitability and does not constitute an application-rate claim.
func MatchFertilizerToNutrient(fertilizerID, fertilizerName, nutrientCode string, percentage, targetAmount float64) (FertilizerMatch, error) {
    nutrientCode = strings.TrimSpace(nutrientCode)
    if nutrientCode == "" { return FertilizerMatch{}, errors.New("nutrient code is required") }
    if percentage <= 0 || percentage > 100 { return FertilizerMatch{}, errors.New("fertilizer percentage must be greater than 0 and at most 100") }
    if targetAmount < 0 { return FertilizerMatch{}, errors.New("target amount cannot be negative") }
    amount, err := CalculateFertilizerAmount(fertilizerID, fertilizerName, nutrientCode, targetAmount, percentage, "kg")
    if err != nil { return FertilizerMatch{}, err }
    return FertilizerMatch{FertilizerID: fertilizerID, FertilizerName: fertilizerName, NutrientCode: nutrientCode, Percentage: percentage, TargetAmount: targetAmount, CalculatedProductAmount: amount.RequiredFertilizerAmount}, nil
}
