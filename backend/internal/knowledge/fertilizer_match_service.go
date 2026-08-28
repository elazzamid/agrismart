package knowledge

import (
    "context"
    "strings"
)

type FertilizerCandidate struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Formulation string `json:"formulation,omitempty"`
    Components []FertilizerComponent `json:"components"`
}

// FindFertilizerCandidates finds products containing at least one requested
// nutrient. It is intentionally a candidate lookup, not an agronomic ranking.
func (s *Service) FindFertilizerCandidates(ctx context.Context, nutrientCodes []string, limit int) ([]FertilizerCandidate, error) {
    if limit <= 0 || limit > 100 { limit = 50 }
    normalized := make([]string, 0, len(nutrientCodes))
    seen := map[string]bool{}
    for _, code := range nutrientCodes {
        code = strings.TrimSpace(code)
        if code != "" && !seen[code] { normalized = append(normalized, code); seen[code] = true }
    }
    if len(normalized) == 0 { return []FertilizerCandidate{}, nil }
    rows, err := s.db.Query(ctx, `
SELECT f.id::text, f.name, COALESCE(f.formulation,''), fn.nutrient_code, COALESCE(fn.percentage,0)
FROM fertilizers f
JOIN fertilizer_nutrients fn ON fn.fertilizer_id=f.id
WHERE fn.nutrient_code = ANY($1)
ORDER BY f.name, fn.nutrient_code
LIMIT $2`, normalized, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]FertilizerCandidate, 0)
    byID := make(map[string]int)
    for rows.Next() {
        var item FertilizerCandidate
        var code string
        var pct float64
        if err := rows.Scan(&item.ID,&item.Name,&item.Formulation,&code,&pct); err != nil { return nil, err }
        idx, ok := byID[item.ID]
        if !ok { result=append(result,item); idx=len(result)-1; byID[item.ID]=idx }
        result[idx].Components=append(result[idx].Components,FertilizerComponent{NutrientCode:code,Percentage:pct})
    }
    if err := rows.Err(); err != nil { return nil, err }
    return result,nil
}
