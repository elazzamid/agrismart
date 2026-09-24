package farm

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Crop struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type CropVariety struct {
	ID          string `json:"id"`
	CropID      string `json:"crop_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
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

type CatalogService struct{ db queryer }

func NewCatalogService(db queryer) *CatalogService { return &CatalogService{db: db} }

func (s *CatalogService) ListCrops(ctx context.Context) ([]Crop, error) {
	rows, err := s.db.Query(ctx, `SELECT id, code, name, COALESCE(description, ''), is_active FROM crops WHERE is_active = TRUE ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	crops := make([]Crop, 0)
	for rows.Next() {
		var crop Crop
		if err := rows.Scan(&crop.ID, &crop.Code, &crop.Name, &crop.Description, &crop.IsActive); err != nil {
			return nil, err
		}
		crops = append(crops, crop)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return crops, nil
}

func (s *CatalogService) GetCrop(ctx context.Context, cropID string) (Crop, error) {
	var crop Crop
	err := s.db.QueryRow(ctx, `SELECT id, code, name, COALESCE(description, ''), is_active FROM crops WHERE id = $1 AND is_active = TRUE`, cropID).
		Scan(&crop.ID, &crop.Code, &crop.Name, &crop.Description, &crop.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return Crop{}, ErrNotFound
	}
	return crop, err
}

func (s *CatalogService) ListVarieties(ctx context.Context, cropID string) ([]CropVariety, error) {
	rows, err := s.db.Query(ctx, `SELECT id, crop_id, name, COALESCE(description, ''), is_active FROM crop_varieties WHERE crop_id = $1 AND is_active = TRUE ORDER BY name`, cropID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	varieties := make([]CropVariety, 0)
	for rows.Next() {
		var variety CropVariety
		if err := rows.Scan(&variety.ID, &variety.CropID, &variety.Name, &variety.Description, &variety.IsActive); err != nil {
			return nil, err
		}
		varieties = append(varieties, variety)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return varieties, nil
}

func (s *CatalogService) ListGrowthStages(ctx context.Context, cropID string) ([]CropGrowthStage, error) {
	rows, err := s.db.Query(ctx, `SELECT id, crop_id, name, sequence_no, min_days, max_days, COALESCE(description, '') FROM crop_growth_stages WHERE crop_id = $1 ORDER BY sequence_no`, cropID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stages := make([]CropGrowthStage, 0)
	for rows.Next() {
		var stage CropGrowthStage
		var minDays, maxDays *int32
		if err := rows.Scan(&stage.ID, &stage.CropID, &stage.Name, &stage.SequenceNo, &minDays, &maxDays, &stage.Description); err != nil {
			return nil, err
		}
		if minDays != nil {
			value := int(*minDays)
			stage.MinDays = &value
		}
		if maxDays != nil {
			value := int(*maxDays)
			stage.MaxDays = &value
		}
		stages = append(stages, stage)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stages, nil
}
