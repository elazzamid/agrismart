package knowledge

import "context"

type FertilizerCandidateResult struct {
	FertilizerCandidate
	Coverages []NutrientCoverage `json:"coverages"`
	SourceDocumentIDs []string `json:"source_document_ids"`
	SourceVersionIDs []string `json:"source_version_ids"`
	AgronomicSuitabilitySupported bool `json:"agronomic_suitability_supported"`
}

// BuildFertilizerCandidates combines published, approved nutrient requirements
// with catalog composition. It intentionally returns candidates and coverage,
// never an inferred agronomic application rate.
func (s *Service) BuildFertilizerCandidates(ctx context.Context, cropID, growthStageID string, productAmount float64, limit int) ([]FertilizerCandidateResult, error) {
	requirements, err := s.ListNutrientRequirements(ctx, cropID, growthStageID, 100)
	if err != nil { return nil, err }
	if len(requirements) == 0 { return []FertilizerCandidateResult{}, nil }
	codes := make([]string,0,len(requirements))
	targets := make([]NutrientTarget,0,len(requirements))
	seen := map[string]bool{}
	for _, r := range requirements {
		if !seen[r.NutrientCode] {
			codes=append(codes,r.NutrientCode); seen[r.NutrientCode]=true
		}
		if r.RequirementMin != nil { targets=append(targets,NutrientTarget{NutrientCode:r.NutrientCode,TargetAmount:*r.RequirementMin,Unit:r.Unit}) }
	}
	if len(targets)==0 { return []FertilizerCandidateResult{}, nil }
	candidates, err := s.FindFertilizerCandidates(ctx, codes, limit)
	if err != nil { return nil, err }
	result:=make([]FertilizerCandidateResult,0,len(candidates))
	for _, c := range candidates {
		analysis, err := AnalyzeMultiNutrientCoverage(c.ID,c.Name,productAmount,c.Components,targets)
		if err != nil { return nil,err }
		result=append(result,FertilizerCandidateResult{FertilizerCandidate:c,Coverages:analysis.Coverages,AgronomicSuitabilitySupported:false})
	}
	return result,nil
}
