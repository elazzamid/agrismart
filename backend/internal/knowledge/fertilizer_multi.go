package knowledge

import "errors"

type NutrientTarget struct { NutrientCode string `json:"nutrient_code"`; TargetAmount float64 `json:"target_amount"`; Unit string `json:"unit"` }
type NutrientCoverage struct { NutrientCode string `json:"nutrient_code"`; TargetAmount float64 `json:"target_amount"`; ProductAmount float64 `json:"product_amount"`; ProvidedAmount float64 `json:"provided_amount"`; CoverageRatio float64 `json:"coverage_ratio"` }
type MultiNutrientMatch struct { FertilizerID string `json:"fertilizer_id"`; FertilizerName string `json:"fertilizer_name"`; ProductAmount float64 `json:"product_amount"`; Coverages []NutrientCoverage `json:"coverages"` }

// AnalyzeMultiNutrientCoverage evaluates a fixed product quantity against
// nutrient targets. It does not infer agronomic suitability or application rate.
func AnalyzeMultiNutrientCoverage(fertilizerID, fertilizerName string, productAmount float64, components []FertilizerComponent, targets []NutrientTarget) (MultiNutrientMatch, error) {
    if productAmount < 0 { return MultiNutrientMatch{}, errors.New("product amount cannot be negative") }
    if len(targets)==0 { return MultiNutrientMatch{}, errors.New("at least one nutrient target is required") }
    seen := make(map[string]bool, len(components))
    for _, c := range components {
        if c.NutrientCode == "" || c.Percentage <= 0 || c.Percentage > 100 { return MultiNutrientMatch{}, errors.New("invalid fertilizer component") }
        if seen[c.NutrientCode] { return MultiNutrientMatch{}, errors.New("duplicate fertilizer nutrient component") }
        seen[c.NutrientCode] = true
    }
    coverage := make([]NutrientCoverage,0,len(targets))
    for _, target := range targets {
        if target.TargetAmount < 0 { return MultiNutrientMatch{}, errors.New("target amount cannot be negative") }
        var pct float64
        for _, c := range components { if c.NutrientCode == target.NutrientCode { pct = c.Percentage; break } }
        if pct <= 0 { coverage = append(coverage, NutrientCoverage{NutrientCode:target.NutrientCode,TargetAmount:target.TargetAmount,ProductAmount:productAmount}); continue }
        provided := productAmount * pct / 100
        ratio := 0.0
        if target.TargetAmount > 0 { ratio = provided / target.TargetAmount }
        coverage = append(coverage, NutrientCoverage{NutrientCode:target.NutrientCode,TargetAmount:target.TargetAmount,ProductAmount:productAmount,ProvidedAmount:provided,CoverageRatio:ratio})
    }
    return MultiNutrientMatch{FertilizerID:fertilizerID,FertilizerName:fertilizerName,ProductAmount:productAmount,Coverages:coverage},nil
}
