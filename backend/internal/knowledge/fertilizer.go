package knowledge

import (
    "errors"
    "math"
    "strings"
)

type FertilizerComponent struct { NutrientCode string `json:"nutrient_code"`; Percentage float64 `json:"percentage"` }
type FertilizerAmount struct { FertilizerID string `json:"fertilizer_id"`; FertilizerName string `json:"fertilizer_name"`; NutrientCode string `json:"nutrient_code"`; TargetAmount float64 `json:"target_amount"`; Unit string `json:"unit"`; RequiredFertilizerAmount float64 `json:"required_fertilizer_amount"` }

// CalculateFertilizerAmount converts a nutrient target into product quantity.
// targetAmount and percentage must use compatible mass units. This is a
// mathematical conversion only and is not an agronomic application-rate claim.
func CalculateFertilizerAmount(fertilizerID, fertilizerName, nutrientCode string, targetAmount, percentage float64, unit string) (FertilizerAmount, error) {
    nutrientCode = strings.TrimSpace(nutrientCode)
    unit = strings.TrimSpace(unit)
    if nutrientCode == "" || unit == "" { return FertilizerAmount{}, errors.New("nutrient code and unit are required") }
    if targetAmount < 0 { return FertilizerAmount{}, errors.New("target amount cannot be negative") }
    if percentage <= 0 || percentage > 100 { return FertilizerAmount{}, errors.New("fertilizer percentage must be greater than 0 and at most 100") }
    amount := targetAmount / (percentage / 100)
    return FertilizerAmount{FertilizerID: fertilizerID, FertilizerName: fertilizerName, NutrientCode: nutrientCode, TargetAmount: targetAmount, Unit: unit, RequiredFertilizerAmount: math.Round(amount*10000)/10000}, nil
}
