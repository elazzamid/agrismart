package farm

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFarmHarvestNotFound   = errors.New("crop cycle not found")
	ErrInvalidFarmHarvestInput = errors.New("invalid farm harvest input")
)

type FarmHarvest struct {
	ID           string  `json:"id"`
	CropCycleID  string  `json:"crop_cycle_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	HarvestDate  string  `json:"harvest_date"`
	QualityGrade string  `json:"quality_grade,omitempty"`
	Notes        string  `json:"notes,omitempty"`
}

type CreateFarmHarvestInput struct {
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	HarvestDate  string  `json:"harvest_date"`
	QualityGrade string  `json:"quality_grade"`
	Notes        string  `json:"notes"`
}

type FarmHarvestService struct{ db queryer }

func NewFarmHarvestService(db queryer) *FarmHarvestService {
	return &FarmHarvestService{db: db}
}

func normalizeFarmHarvestInput(in CreateFarmHarvestInput) (CreateFarmHarvestInput, error) {
	in.Unit = strings.TrimSpace(in.Unit)
	in.HarvestDate = strings.TrimSpace(in.HarvestDate)
	in.QualityGrade = strings.TrimSpace(in.QualityGrade)
	in.Notes = strings.TrimSpace(in.Notes)

	if math.IsNaN(in.Quantity) || math.IsInf(in.Quantity, 0) || in.Quantity < 0.001 || in.Quantity > 99999999999.999 {
		return CreateFarmHarvestInput{}, ErrInvalidFarmHarvestInput
	}
	if utf8.RuneCountInString(in.Unit) < 1 || utf8.RuneCountInString(in.Unit) > 20 {
		return CreateFarmHarvestInput{}, ErrInvalidFarmHarvestInput
	}
	parsed, err := time.Parse("2006-01-02", in.HarvestDate)
	if err != nil || parsed.Format("2006-01-02") != in.HarvestDate {
		return CreateFarmHarvestInput{}, ErrInvalidFarmHarvestInput
	}
	if utf8.RuneCountInString(in.QualityGrade) > 40 || (in.QualityGrade != "" && utf8.RuneCountInString(in.QualityGrade) < 1) {
		return CreateFarmHarvestInput{}, ErrInvalidFarmHarvestInput
	}
	if utf8.RuneCountInString(in.Notes) > 2000 {
		return CreateFarmHarvestInput{}, ErrInvalidFarmHarvestInput
	}
	return in, nil
}

// Create records an actual harvest for a crop cycle owned by the authenticated farmer.
func (s *FarmHarvestService) Create(ctx context.Context, userID, cropCycleID string, in CreateFarmHarvestInput) (FarmHarvest, error) {
	normalized, err := normalizeFarmHarvestInput(in)
	if err != nil {
		return FarmHarvest{}, err
	}

	var harvest FarmHarvest
	err = s.db.QueryRow(ctx, `
		INSERT INTO farm_harvests (crop_cycle_id, quantity, unit, harvest_date, quality_grade, notes)
		SELECT cc.id, $3, $4, $5::date, NULLIF($6, ''), NULLIF($7, '')
		FROM crop_cycles cc
		JOIN farm_plots p ON p.id = cc.plot_id
		JOIN farms f ON f.id = p.farm_id
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE cc.id = $1 AND fp.user_id = $2
		RETURNING id, crop_cycle_id, quantity, unit, harvest_date::text, COALESCE(quality_grade, ''), COALESCE(notes, '')`,
		cropCycleID, userID, normalized.Quantity, normalized.Unit, normalized.HarvestDate,
		normalized.QualityGrade, normalized.Notes,
	).Scan(&harvest.ID, &harvest.CropCycleID, &harvest.Quantity, &harvest.Unit,
		&harvest.HarvestDate, &harvest.QualityGrade, &harvest.Notes)
	if errors.Is(err, pgx.ErrNoRows) {
		return FarmHarvest{}, ErrFarmHarvestNotFound
	}
	return harvest, err
}

func (s *FarmHarvestService) List(ctx context.Context, userID, cropCycleID string) ([]FarmHarvest, error) {
	var owned bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM crop_cycles cc
			JOIN farm_plots p ON p.id = cc.plot_id
			JOIN farms f ON f.id = p.farm_id
			JOIN farmer_profiles fp ON fp.id = f.farmer_id
			WHERE cc.id = $1 AND fp.user_id = $2
		)`, cropCycleID, userID).Scan(&owned)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, ErrFarmHarvestNotFound
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, crop_cycle_id, quantity, unit, harvest_date::text, COALESCE(quality_grade, ''), COALESCE(notes, '')
		FROM farm_harvests
		WHERE crop_cycle_id = $1
		ORDER BY harvest_date DESC, created_at DESC`, cropCycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	harvests := make([]FarmHarvest, 0)
	for rows.Next() {
		var harvest FarmHarvest
		if err := rows.Scan(&harvest.ID, &harvest.CropCycleID, &harvest.Quantity, &harvest.Unit,
			&harvest.HarvestDate, &harvest.QualityGrade, &harvest.Notes); err != nil {
			return nil, err
		}
		harvests = append(harvests, harvest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return harvests, nil
}
