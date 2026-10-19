package farm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCropCycleNotFound = errors.New("crop cycle not found")
	ErrInvalidCropCycleInput = errors.New("invalid crop cycle input")
)

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

type CropCycleService struct{ db queryer }

func NewCropCycleService(db queryer) *CropCycleService { return &CropCycleService{db: db} }

// Create scopes the operation using the authenticated users.id. The database
// relationship is resolved through farmer_profiles.user_id.
func (s *CropCycleService) Create(ctx context.Context, userID, plotID string, in CreateCropCycleInput) (CropCycle, error) {
	if strings.TrimSpace(in.CropID) == "" {
		return CropCycle{}, ErrInvalidCropCycleInput
	}
	if _, err := time.Parse("2006-01-02", in.PlantingDate); err != nil {
		return CropCycle{}, ErrInvalidCropCycleInput
	}
	var c CropCycle
	err := s.db.QueryRow(ctx, `
		INSERT INTO crop_cycles (plot_id, crop_id, variety_id, planting_date)
		SELECT p.id, $3, $4, $5::date
		FROM farm_plots p
		JOIN farms f ON f.id = p.farm_id
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE p.id = $1 AND fp.user_id = $2
		RETURNING id, plot_id, crop_id, variety_id, planting_date::text, status`,
		plotID, userID, in.CropID, in.VarietyID, in.PlantingDate,
	).Scan(&c.ID, &c.PlotID, &c.CropID, &c.VarietyID, &c.PlantingDate, &c.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return CropCycle{}, ErrCropCycleNotFound
	}
	return c, err
}

func (s *CropCycleService) List(ctx context.Context, userID, plotID string) ([]CropCycle, error) {
	rows, err := s.db.Query(ctx, `
		SELECT cc.id, cc.plot_id, cc.crop_id, cc.variety_id, cc.planting_date::text, cc.status
		FROM crop_cycles cc
		JOIN farm_plots p ON p.id = cc.plot_id
		JOIN farms f ON f.id = p.farm_id
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE cc.plot_id = $1 AND fp.user_id = $2
		ORDER BY cc.planting_date DESC`, plotID, userID)
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
