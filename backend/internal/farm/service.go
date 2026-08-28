package farm

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("farm not found")

// Farm is the farmer-owned agricultural location.
type Farm struct {
	ID           string  `json:"id"`
	FarmerID     string  `json:"farmer_id"`
	Name         string  `json:"name"`
	LocationName string  `json:"location_name,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
}

type CreateInput struct {
	Name         string  `json:"name"`
	LocationName string  `json:"location_name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type Service struct { db *pgxpool.Pool }

func NewService(db *pgxpool.Pool) *Service { return &Service{db: db} }

func (s *Service) Create(ctx context.Context, farmerID string, in CreateInput) (Farm, error) {
	if strings.TrimSpace(in.Name) == "" {
		return Farm{}, errors.New("farm name is required")
	}
	var f Farm
	err := s.db.QueryRow(ctx, `
		INSERT INTO farms (farmer_id, name, location_name, latitude, longitude)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5)
		RETURNING id, farmer_id, name, COALESCE(location_name, ''), latitude, longitude`,
		farmerID, strings.TrimSpace(in.Name), strings.TrimSpace(in.LocationName), in.Latitude, in.Longitude,
	).Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude)
	return f, err
}

func (s *Service) List(ctx context.Context, farmerID string) ([]Farm, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, farmer_id, name, COALESCE(location_name, ''), latitude, longitude
		FROM farms WHERE farmer_id = $1 ORDER BY created_at DESC`, farmerID)
	if err != nil { return nil, err }
	defer rows.Close()
	farms := make([]Farm, 0)
	for rows.Next() {
		var f Farm
		if err := rows.Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude); err != nil { return nil, err }
		farms = append(farms, f)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return farms, nil
}

func (s *Service) Get(ctx context.Context, farmerID, id string) (Farm, error) {
	var f Farm
	err := s.db.QueryRow(ctx, `
		SELECT id, farmer_id, name, COALESCE(location_name, ''), latitude, longitude
		FROM farms WHERE id = $1 AND farmer_id = $2`, id, farmerID).
		Scan(&f.ID, &f.FarmerID, &f.Name, &f.LocationName, &f.Latitude, &f.Longitude)
	if errors.Is(err, pgx.ErrNoRows) { return Farm{}, ErrNotFound }
	return f, err
}
