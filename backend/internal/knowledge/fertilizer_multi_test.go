package knowledge

import "testing"

func TestAnalyzeMultiNutrientCoverage(t *testing.T) {
    got, err := AnalyzeMultiNutrientCoverage("f1", "NPK 15-15-15", 100, []FertilizerComponent{{NutrientCode:"N",Percentage:15},{NutrientCode:"P",Percentage:15},{NutrientCode:"K",Percentage:15}}, []NutrientTarget{{NutrientCode:"N",TargetAmount:10,Unit:"kg"},{NutrientCode:"P",TargetAmount:20,Unit:"kg"},{NutrientCode:"K",TargetAmount:5,Unit:"kg"}})
    if err != nil { t.Fatal(err) }
    if len(got.Coverages) != 3 { t.Fatalf("expected 3 coverages, got %d",len(got.Coverages)) }
    if got.Coverages[0].ProvidedAmount != 15 { t.Fatalf("unexpected N provided: %v",got.Coverages[0].ProvidedAmount) }
    if got.Coverages[1].CoverageRatio != 0.75 { t.Fatalf("unexpected P ratio: %v",got.Coverages[1].CoverageRatio) }
}
