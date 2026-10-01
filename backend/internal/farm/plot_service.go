package farm

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

var ErrPlotNotFound = errors.New("plot not found")

type Plot struct {
	ID     string  `json:"id"`
	FarmID string  `json:"farm_id"`
	Name   string  `json:"name"`
	AreaM2 float64 `json:"area_m2"`
}

type CreatePlotInput struct {
	Name   string  `json:"name"`
	AreaM2 float64 `json:"area_m2"`
}

type PlotService struct { db DBTX }
func NewPlotService(db DBTX) *PlotService { return &PlotService{db: db} }

func (s *PlotService) Create(ctx context.Context, farmerID, farmID string, in CreatePlotInput) (Plot, error) {
	if strings.TrimSpace(in.Name) == "" { return Plot{}, errors.New("plot name is required") }
	if in.AreaM2 <= 0 { return Plot{}, errors.New("plot area must be greater than zero") }
	var p Plot
	err := s.db.QueryRow(ctx, `
		INSERT INTO farm_plots (farm_id, name, area_m2)
		SELECT id, $3, $4 FROM farms WHERE id = $1 AND farmer_id = $2
		RETURNING id, farm_id, name, area_m2`, farmID, farmerID, strings.TrimSpace(in.Name), in.AreaM2).
		Scan(&p.ID, &p.FarmID, &p.Name, &p.AreaM2)
	if errors.Is(err, pgx.ErrNoRows) { return Plot{}, ErrPlotNotFound }
	return p, err
}

func (s *PlotService) List(ctx context.Context, farmerID, farmID string) ([]Plot, error) {
	rows, err := s.db.Query(ctx, `
		SELECT p.id, p.farm_id, p.name, p.area_m2
		FROM farm_plots p JOIN farms f ON f.id = p.farm_id
		WHERE p.farm_id = $1 AND f.farmer_id = $2 ORDER BY p.created_at DESC`, farmID, farmerID)
	if err != nil { return nil, err }
	defer rows.Close()
	plots := make([]Plot, 0)
	for rows.Next() {
		var p Plot
		if err := rows.Scan(&p.ID, &p.FarmID, &p.Name, &p.AreaM2); err != nil { return nil, err }
		plots = append(plots, p)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return plots, nil
}
