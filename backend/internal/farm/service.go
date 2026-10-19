package farm

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a farm is not owned by the authenticated user.
var ErrNotFound = errors.New("farm not found")

// queryer captures the database operations used by the farm services and keeps
// the services directly testable with pgxmock.
type queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Farm is the farmer-owned agricultural location.
type Farm struct {
	ID           string   `json:"id"`
	FarmerID     string   `json:"farmer_id"`
	Name         string   `json:"name"`
	LocationName string   `json:"location_name,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
}

type CreateInput struct {
	Name         string   `json:"name"`
	LocationName string   `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type Service struct{ db queryer }

func NewService(db queryer) *Service { return &Service{db: db} }

// Create resolves the authenticated users.id to farmer_profiles.id before
// writing farms.farmer_id, which references farmer_profiles rather than users.
func (s *Service) Create(ctx context.Context, userID string, in CreateInput) (Farm, error) {
	if strings.TrimSpace(in.Name) == "" {
		return Farm{}, errors.New("farm name is required")
	}
	var f Farm
	err := s.db.QueryRow(ctx, `
		INSERT INTO farms (farmer_id, name, location_name, latitude, longitude)
		SELECT fp.id, $2, NULLIF($3, ''), $4, $5
		FROM farmer_profiles fp
		WHERE fp.user_id = $1
		RETURNING id, farmer_id, name, COALESCE(location_name, ''), latitude, longitude`,
		userID, strings.TrimSpace(in.Name), strings.TrimSpace(in.LocationName), in.Latitude, in.Longitude,
	).Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude)
	return f, err
}

func (s *Service) List(ctx context.Context, userID string) ([]Farm, error) {
	rows, err := s.db.Query(ctx, `
		SELECT f.id, f.farmer_id, f.name, COALESCE(f.location_name, ''), f.latitude, f.longitude
		FROM farms f
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE fp.user_id = $1
		ORDER BY f.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	farms := make([]Farm, 0)
	for rows.Next() {
		var f Farm
		if err := rows.Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude); err != nil {
			return nil, err
		}
		farms = append(farms, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return farms, nil
}

func (s *Service) Get(ctx context.Context, userID, id string) (Farm, error) {
	var f Farm
	err := s.db.QueryRow(ctx, `
		SELECT f.id, f.farmer_id, f.name, COALESCE(f.location_name, ''), f.latitude, f.longitude
		FROM farms f
		JOIN farmer_profiles fp ON fp.id = f.farmer_id
		WHERE f.id = $1 AND fp.user_id = $2`, id, userID).
		Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude)
	if errors.Is(err, pgx.ErrNoRows) {
		return Farm{}, ErrNotFound
	}
	return f, err
}
