package farm

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFarmActivityNotFound = errors.New("crop cycle not found")
	ErrInvalidFarmActivityInput = errors.New("invalid farm activity input")
)

type FarmActivity struct {
	ID string `json:"id"`
	CropCycleID string `json:"crop_cycle_id"`
	ActivityType string `json:"activity_type"`
	ActivityDate string `json:"activity_date"`
	Notes string `json:"notes,omitempty"`
}

type CreateFarmActivityInput struct {
	ActivityType string `json:"activity_type"`
	ActivityDate string `json:"activity_date"`
	Notes string `json:"notes"`
}

type FarmActivityService struct { db queryer }
func NewFarmActivityService(db queryer) *FarmActivityService { return &FarmActivityService{db: db} }

func normalizeFarmActivityInput(in CreateFarmActivityInput) (CreateFarmActivityInput, error) {
	in.ActivityType = strings.ToLower(strings.TrimSpace(in.ActivityType))
	in.ActivityDate = strings.TrimSpace(in.ActivityDate)
	in.Notes = strings.TrimSpace(in.Notes)
	if utf8.RuneCountInString(in.ActivityType) < 1 || utf8.RuneCountInString(in.ActivityType) > 80 { return CreateFarmActivityInput{}, ErrInvalidFarmActivityInput }
	if parsed, err := time.Parse("2006-01-02", in.ActivityDate); err != nil || parsed.Format("2006-01-02") != in.ActivityDate { return CreateFarmActivityInput{}, ErrInvalidFarmActivityInput }
	if utf8.RuneCountInString(in.Notes) > 2000 { return CreateFarmActivityInput{}, ErrInvalidFarmActivityInput }
	return in, nil
}

func (s *FarmActivityService) Create(ctx context.Context, userID, cropCycleID string, in CreateFarmActivityInput) (FarmActivity, error) {
	normalized, err := normalizeFarmActivityInput(in)
	if err != nil { return FarmActivity{}, err }
	var activity FarmActivity
	err = s.db.QueryRow(ctx, `INSERT INTO farm_activities (crop_cycle_id, activity_type, activity_date, notes) SELECT cc.id, $3, $4::date, NULLIF($5, '') FROM crop_cycles cc JOIN farm_plots p ON p.id = cc.plot_id JOIN farms f ON f.id = p.farm_id JOIN farmer_profiles fp ON fp.id = f.farmer_id WHERE cc.id = $1 AND fp.user_id = $2 RETURNING id, crop_cycle_id, activity_type, activity_date::text, COALESCE(notes, '')`, cropCycleID, userID, normalized.ActivityType, normalized.ActivityDate, normalized.Notes).Scan(&activity.ID, &activity.CropCycleID, &activity.ActivityType, &activity.ActivityDate, &activity.Notes)
	if errors.Is(err, pgx.ErrNoRows) { return FarmActivity{}, ErrFarmActivityNotFound }
	return activity, err
}

func (s *FarmActivityService) List(ctx context.Context, userID, cropCycleID string) ([]FarmActivity, error) {
	var owned bool
	err := s.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM crop_cycles cc JOIN farm_plots p ON p.id = cc.plot_id JOIN farms f ON f.id = p.farm_id JOIN farmer_profiles fp ON fp.id = f.farmer_id WHERE cc.id = $1 AND fp.user_id = $2)`, cropCycleID, userID).Scan(&owned)
	if err != nil { return nil, err }
	if !owned { return nil, ErrFarmActivityNotFound }
	rows, err := s.db.Query(ctx, `SELECT id, crop_cycle_id, activity_type, activity_date::text, COALESCE(notes, '') FROM farm_activities WHERE crop_cycle_id = $1 ORDER BY activity_date DESC, created_at DESC`, cropCycleID)
	if err != nil { return nil, err }
	defer rows.Close()
	activities := make([]FarmActivity, 0)
	for rows.Next() { var activity FarmActivity; if err := rows.Scan(&activity.ID, &activity.CropCycleID, &activity.ActivityType, &activity.ActivityDate, &activity.Notes); err != nil { return nil, err }; activities = append(activities, activity) }
	if err := rows.Err(); err != nil { return nil, err }
	return activities, nil
}
