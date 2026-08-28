package knowledge

import "context"

type FertilizerCatalogItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Formulation string `json:"formulation,omitempty"`
	Description string `json:"description,omitempty"`
	Components []FertilizerComponent `json:"components"`
}

// ListFertilizers returns fertilizer products and their nutrient composition.
// This is catalog data only; it does not rank agronomic suitability.
func (s *Service) ListFertilizers(ctx context.Context, limit int) ([]FertilizerCatalogItem, error) {
	if limit <= 0 || limit > 100 { limit = 50 }
	rows, err := s.db.Query(ctx, `
SELECT f.id::text, f.name, COALESCE(f.formulation,''), COALESCE(f.description,''),
       COALESCE(fn.nutrient_code,''), COALESCE(fn.percentage,0)
FROM fertilizers f
LEFT JOIN fertilizer_nutrients fn ON fn.fertilizer_id=f.id
ORDER BY f.name, fn.nutrient_code
LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	byID := make(map[string]int)
	result := make([]FertilizerCatalogItem, 0)
	for rows.Next() {
		var item FertilizerCatalogItem
		var code string
		var percentage float64
		if err := rows.Scan(&item.ID,&item.Name,&item.Formulation,&item.Description,&code,&percentage); err != nil { return nil, err }
		idx, ok := byID[item.ID]
		if !ok { result=append(result,item); idx=len(result)-1; byID[item.ID]=idx }
		if code != "" { result[idx].Components=append(result[idx].Components,FertilizerComponent{NutrientCode:code,Percentage:percentage}) }
	}
	if err := rows.Err(); err != nil { return nil, err }
	return result,nil
}
