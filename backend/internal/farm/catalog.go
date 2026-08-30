package farm

import "context"

type Crop struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CropVariety struct {
	ID          string `json:"id"`
	CropID      string `json:"crop_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CropGrowthStage struct {
	ID          string `json:"id"`
	CropID      string `json:"crop_id"`
	Name        string `json:"name"`
	SequenceNo  int    `json:"sequence_no"`
	MinDays     *int   `json:"min_days,omitempty"`
	MaxDays     *int   `json:"max_days,omitempty"`
	Description string `json:"description,omitempty"`
}

type CatalogService struct{ db DBTX }

func NewCatalogService(db DBTX) *CatalogService { return &CatalogService{db: db} }

func (s *CatalogService) ListCrops(ctx context.Context) ([]Crop, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, code, name, COALESCE(description, '')
		FROM crops WHERE is_active = TRUE ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]Crop, 0)
	for rows.Next() {
		var item Crop
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Description); err != nil { return nil, err }
		items = append(items, item)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return items, nil
}

func (s *CatalogService) ListVarieties(ctx context.Context, cropID string) ([]CropVariety, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, crop_id, name, COALESCE(description, '')
		FROM crop_varieties
		WHERE crop_id = $1 AND is_active = TRUE ORDER BY name`, cropID)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]CropVariety, 0)
	for rows.Next() {
		var item CropVariety
		if err := rows.Scan(&item.ID, &item.CropID, &item.Name, &item.Description); err != nil { return nil, err }
		items = append(items, item)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return items, nil
}

func (s *CatalogService) ListGrowthStages(ctx context.Context, cropID string) ([]CropGrowthStage, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, crop_id, name, sequence_no, min_days, max_days, COALESCE(description, '')
		FROM crop_growth_stages
		WHERE crop_id = $1 ORDER BY sequence_no`, cropID)
	if err != nil { return nil, err }
	defer rows.Close()
	items := make([]CropGrowthStage, 0)
	for rows.Next() {
		var item CropGrowthStage
		if err := rows.Scan(&item.ID, &item.CropID, &item.Name, &item.SequenceNo, &item.MinDays, &item.MaxDays, &item.Description); err != nil { return nil, err }
		items = append(items, item)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return items, nil
}
