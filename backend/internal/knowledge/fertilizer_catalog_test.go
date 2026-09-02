package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestListFertilizersLimitAppliesToProductsBeforeNutrientJoin(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	s := NewService(pool)
	pool.ExpectQuery(`WITH limited_fertilizers AS \(\s*SELECT id, name, formulation, description\s*FROM fertilizers\s*ORDER BY name, id\s*LIMIT \$1\s*\)`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "formulation", "description", "nutrient_code", "percentage"}).
			AddRow("fert-1", "Urea", "46-0-0", "", "N", 46.0).
			AddRow("fert-1", "Urea", "46-0-0", "", "", 0.0))

	items, err := s.ListFertilizers(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one fertilizer, got %d", len(items))
	}
	if items[0].ID != "fert-1" || len(items[0].Components) != 1 {
		t.Fatalf("unexpected catalog item: %+v", items[0])
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
