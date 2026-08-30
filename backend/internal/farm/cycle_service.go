package farm

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

var ErrCropCycleNotFound = errors.New("crop cycle not found")

type CropCycle struct {
	ID           string  `json:"id"`
	PlotID       string  `json:"plot_id"`
	CropID       string  `json:"crop_id"`
	VarietyID    *string `json:"variety_id,omitempty"`
	PlantingDate string  `json:"planting_date"`
	Status       string  `json:"status"`
}

type CreateCropCycleInput struct {
	CropID       string  `json:"crop_id"`
	VarietyID    *string `json:"variety_id"`
	PlantingDate string  `json:"planting_date"`
}

type CropCycleService struct{ db DBTX }

func NewCropCycleService(db DBTX) *CropCycleService { return &CropCycleService{db: db} }

func (s *CropCycleService) Create(ctx context.Context, farmerID, plotID string, in CreateCropCycleInput) (CropCycle, error) {
	if strings.TrimSpace(in.CropID) == "" {
		return CropCycle{}, errors.New("crop id is required")
	}
	if strings.TrimSpace(in.PlantingDate) == "" {
		return CropCycle{}, errors.New("planting date is required")
	}

	var varietyID any
	if in.VarietyID != nil {
		varietyID = *in.VarietyID
	}

	var c CropCycle
	err := s.db.QueryRow(ctx, `
		INSERT INTO crop_cycles (plot_id, crop_id, variety_id, planting_date)
		SELECT p.id, $3, $4, $5::date
		FROM farm_plots p JOIN farms f ON f.id = p.farm_id
		WHERE p.id = $1 AND f.farmer_id = $2
		RETURNING id, plot_id, crop_id, variety_id, planting_date::text, status`,
		plotID, farmerID, in.CropID, varietyID, in.PlantingDate,
	).Scan(&c.ID, &c.PlotID, &c.CropID, &c.VarietyID, &c.PlantingDate, &c.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return CropCycle{}, ErrCropCycleNotFound
	}
	return c, err
}

func (s *CropCycleService) List(ctx context.Context, farmerID, plotID string) ([]CropCycle, error) {
	rows, err := s.db.Query(ctx, `
		SELECT cc.id, cc.plot_id, cc.crop_id, cc.variety_id, cc.planting_date::text, cc.status
		FROM crop_cycles cc JOIN farm_plots p ON p.id = cc.plot_id
		JOIN farms f ON f.id = p.farm_id
		WHERE cc.plot_id = $1 AND f.farmer_id = $2 ORDER BY cc.planting_date DESC`, plotID, farmerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cycles := make([]CropCycle, 0)
	for rows.Next() {
		var c CropCycle
		if err := rows.Scan(&c.ID, &c.PlotID, &c.CropID, &c.VarietyID, &c.PlantingDate, &c.Status); err != nil {
			return nil, err
		}
		cycles = append(cycles, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cycles, nil
}
