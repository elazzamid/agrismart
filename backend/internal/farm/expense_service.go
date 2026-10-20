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
	ErrFarmExpenseNotFound   = errors.New("crop cycle not found")
	ErrInvalidFarmExpenseInput = errors.New("invalid farm expense input")
)

type FarmExpense struct {
	ID          string  `json:"id"`
	CropCycleID string  `json:"crop_cycle_id"`
	ExpenseType string  `json:"expense_type"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	ExpenseDate string  `json:"expense_date"`
	Notes       string  `json:"notes,omitempty"`
}

type CreateFarmExpenseInput struct {
	ExpenseType string  `json:"expense_type"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	ExpenseDate string  `json:"expense_date"`
	Notes       string  `json:"notes"`
}

type FarmExpenseService struct{ db queryer }

func NewFarmExpenseService(db queryer) *FarmExpenseService {
	return &FarmExpenseService{db: db}
}

func normalizeFarmExpenseInput(in CreateFarmExpenseInput) (CreateFarmExpenseInput, error) {
	in.ExpenseType = strings.ToLower(strings.TrimSpace(in.ExpenseType))
	in.Description = strings.TrimSpace(in.Description)
	in.ExpenseDate = strings.TrimSpace(in.ExpenseDate)
	in.Notes = strings.TrimSpace(in.Notes)

	switch in.ExpenseType {
	case "seed", "fertilizer", "pesticide", "labor", "irrigation", "equipment", "transport", "other":
	default:
		return CreateFarmExpenseInput{}, ErrInvalidFarmExpenseInput
	}
	if utf8.RuneCountInString(in.Description) < 1 || utf8.RuneCountInString(in.Description) > 160 {
		return CreateFarmExpenseInput{}, ErrInvalidFarmExpenseInput
	}
	if math.IsNaN(in.Amount) || math.IsInf(in.Amount, 0) || in.Amount < 0.01 || in.Amount > 999999999999.99 {
		return CreateFarmExpenseInput{}, ErrInvalidFarmExpenseInput
	}
	parsed, err := time.Parse("2006-01-02", in.ExpenseDate)
	if err != nil || parsed.Format("2006-01-02") != in.ExpenseDate {
		return CreateFarmExpenseInput{}, ErrInvalidFarmExpenseInput
	}
	if utf8.RuneCountInString(in.Notes) > 2000 {
		return CreateFarmExpenseInput{}, ErrInvalidFarmExpenseInput
	}
	return in, nil
}

// Create records an expense for a crop cycle owned by the authenticated farmer.
func (s *FarmExpenseService) Create(ctx context.Context, userID, cropCycleID string, in CreateFarmExpenseInput) (FarmExpense, error) {
	normalized, err := normalizeFarmExpenseInput(in)
	if err != nil {
		return FarmExpense{}, err
	}

	var expense FarmExpense
	err = s.db.QueryRow(ctx, `
		INSERT INTO farm_expenses (crop_cycle_id, expense_type, description, amount, expense_date, notes)
		SELECT cc.id, $3, $4, $5, $6::date, NULLIF($7, '')
		FROM crop_cycles cc
		JOIN farm_plots p ON p.id = cc.plot_id
		JOIN farms f ON f.id = p.farm_id
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE cc.id = $1 AND fp.user_id = $2
		RETURNING id, crop_cycle_id, expense_type, description, amount, expense_date::text, COALESCE(notes, '')`,
		cropCycleID, userID, normalized.ExpenseType, normalized.Description,
		normalized.Amount, normalized.ExpenseDate, normalized.Notes,
	).Scan(&expense.ID, &expense.CropCycleID, &expense.ExpenseType, &expense.Description,
		&expense.Amount, &expense.ExpenseDate, &expense.Notes)
	if errors.Is(err, pgx.ErrNoRows) {
		return FarmExpense{}, ErrFarmExpenseNotFound
	}
	return expense, err
}

func (s *FarmExpenseService) List(ctx context.Context, userID, cropCycleID string) ([]FarmExpense, error) {
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
		return nil, ErrFarmExpenseNotFound
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, crop_cycle_id, expense_type, description, amount, expense_date::text, COALESCE(notes, '')
		FROM farm_expenses
		WHERE crop_cycle_id = $1
		ORDER BY expense_date DESC, created_at DESC`, cropCycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]FarmExpense, 0)
	for rows.Next() {
		var expense FarmExpense
		if err := rows.Scan(&expense.ID, &expense.CropCycleID, &expense.ExpenseType, &expense.Description,
			&expense.Amount, &expense.ExpenseDate, &expense.Notes); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}
