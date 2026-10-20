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
	ErrFarmInputNotFound   = errors.New("crop cycle not found")
	ErrInvalidFarmInput    = errors.New("invalid farm input")
)

type FarmInput struct {
	ID          string  `json:"id"`
	CropCycleID string  `json:"crop_cycle_id"`
	InputType   string  `json:"input_type"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	AppliedDate string  `json:"applied_date"`
	Notes       string  `json:"notes,omitempty"`
}

type CreateFarmInputInput struct {
	InputType   string  `json:"input_type"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	AppliedDate string  `json:"applied_date"`
	Notes       string  `json:"notes"`
}

type FarmInputService struct{ db queryer }

func NewFarmInputService(db queryer) *FarmInputService {
	return &FarmInputService{db: db}
}

func normalizeFarmInput(in CreateFarmInputInput) (CreateFarmInputInput, error) {
	in.InputType = strings.ToLower(strings.TrimSpace(in.InputType))
	in.ProductName = strings.TrimSpace(in.ProductName)
	in.Unit = strings.TrimSpace(in.Unit)
	in.AppliedDate = strings.TrimSpace(in.AppliedDate)
	in.Notes = strings.TrimSpace(in.Notes)

	switch in.InputType {
	case "fertilizer", "pesticide", "seed", "other":
	default:
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	if utf8.RuneCountInString(in.ProductName) < 1 || utf8.RuneCountInString(in.ProductName) > 120 {
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	if math.IsNaN(in.Quantity) || math.IsInf(in.Quantity, 0) || in.Quantity < 0.001 {
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	if utf8.RuneCountInString(in.Unit) < 1 || utf8.RuneCountInString(in.Unit) > 20 {
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	parsed, err := time.Parse("2006-01-02", in.AppliedDate)
	if err != nil || parsed.Format("2006-01-02") != in.AppliedDate {
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	if utf8.RuneCountInString(in.Notes) > 2000 {
		return CreateFarmInputInput{}, ErrInvalidFarmInput
	}
	return in, nil
}

// Create records an input actually applied to a crop cycle owned by the caller.
// It does not calculate or recommend an agronomic dose.
func (s *FarmInputService) Create(ctx context.Context, userID, cropCycleID string, in CreateFarmInputInput) (FarmInput, error) {
	normalized, err := normalizeFarmInput(in)
	if err != nil {
		return FarmInput{}, err
	}

	var input FarmInput
	err = s.db.QueryRow(ctx, `
		INSERT INTO farm_inputs (crop_cycle_id, input_type, product_name, quantity, unit, applied_date, notes)
		SELECT cc.id, $3, $4, $5, $6, $7::date, NULLIF($8, '')
		FROM crop_cycles cc
		JOIN farm_plots p ON p.id = cc.plot_id
		JOIN farms f ON f.id = p.farm_id
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE cc.id = $1 AND fp.user_id = $2
		RETURNING id, crop_cycle_id, input_type, product_name, quantity, unit, applied_date::text, COALESCE(notes, '')`,
		cropCycleID, userID, normalized.InputType, normalized.ProductName,
		normalized.Quantity, normalized.Unit, normalized.AppliedDate, normalized.Notes,
	).Scan(&input.ID, &input.CropCycleID, &input.InputType, &input.ProductName,
		&input.Quantity, &input.Unit, &input.AppliedDate, &input.Notes)
	if errors.Is(err, pgx.ErrNoRows) {
		return FarmInput{}, ErrFarmInputNotFound
	}
	return input, err
}

func (s *FarmInputService) List(ctx context.Context, userID, cropCycleID string) ([]FarmInput, error) {
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
		return nil, ErrFarmInputNotFound
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, crop_cycle_id, input_type, product_name, quantity, unit, applied_date::text, COALESCE(notes, '')
		FROM farm_inputs
		WHERE crop_cycle_id = $1
		ORDER BY applied_date DESC, created_at DESC`, cropCycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inputs := make([]FarmInput, 0)
	for rows.Next() {
		var input FarmInput
		if err := rows.Scan(&input.ID, &input.CropCycleID, &input.InputType, &input.ProductName,
			&input.Quantity, &input.Unit, &input.AppliedDate, &input.Notes); err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return inputs, nil
}
