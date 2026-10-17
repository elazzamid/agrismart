package knowledge

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestFindFertilizerCandidatesLimitAppliesToProductsBeforeNutrientJoin(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	s := NewService(pool)
	pool.ExpectQuery(`WITH limited_fertilizers AS \(\s*SELECT f\.id, f\.name, f\.formulation\s*FROM fertilizers f\s*JOIN fertilizer_nutrients fn ON fn\.fertilizer_id=f\.id\s*WHERE fn\.nutrient_code = ANY\(\$1\)\s*GROUP BY f\.id, f\.name, f\.formulation\s*ORDER BY f\.name, f\.id\s*LIMIT \$2\s*\)`).
		WithArgs([]string{"N"}, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "formulation", "nutrient_code", "percentage"}).
			AddRow("fert-1", "Urea", "46-0-0", "N", 46.0))

	items, err := s.FindFertilizerCandidates(context.Background(), []string{" N ", "N"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one fertilizer, got %d", len(items))
	}
	if items[0].ID != "fert-1" || len(items[0].Components) != 1 {
		t.Fatalf("unexpected candidate: %+v", items[0])
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
