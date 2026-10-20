package farm

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestNormalizeFarmHarvestInput(t *testing.T) {
	got, err := normalizeFarmHarvestInput(CreateFarmHarvestInput{
		Quantity: 125.5, Unit: " kg ", HarvestDate: "2026-10-19",
		QualityGrade: " Grade A ", Notes: " First pick ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Unit != "kg" || got.QualityGrade != "Grade A" || got.Notes != "First pick" {
		t.Fatalf("unexpected normalized harvest: %+v", got)
	}
}

func TestNormalizeFarmHarvestInputRejectsInvalidValues(t *testing.T) {
	base := CreateFarmHarvestInput{Quantity: 1, Unit: "kg", HarvestDate: "2026-10-19"}
	tests := []struct {
		name string
		mutate func(*CreateFarmHarvestInput)
	}{
		{name: "zero quantity", mutate: func(in *CreateFarmHarvestInput) { in.Quantity = 0 }},
		{name: "negative quantity", mutate: func(in *CreateFarmHarvestInput) { in.Quantity = -1 }},
		{name: "NaN quantity", mutate: func(in *CreateFarmHarvestInput) { in.Quantity = math.NaN() }},
		{name: "infinite quantity", mutate: func(in *CreateFarmHarvestInput) { in.Quantity = math.Inf(1) }},
		{name: "quantity too large", mutate: func(in *CreateFarmHarvestInput) { in.Quantity = 100000000000 }},
		{name: "empty unit", mutate: func(in *CreateFarmHarvestInput) { in.Unit = " " }},
		{name: "invalid date", mutate: func(in *CreateFarmHarvestInput) { in.HarvestDate = "2026-02-30" }},
		{name: "noncanonical date", mutate: func(in *CreateFarmHarvestInput) { in.HarvestDate = "2026-1-9" }},
		{name: "grade too long", mutate: func(in *CreateFarmHarvestInput) { in.QualityGrade = strings.Repeat("x", 41) }},
		{name: "notes too long", mutate: func(in *CreateFarmHarvestInput) { in.Notes = strings.Repeat("x", 2001) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := base
			tt.mutate(&in)
			if _, err := normalizeFarmHarvestInput(in); err != ErrInvalidFarmHarvestInput {
				t.Fatalf("normalizeFarmHarvestInput() error = %v, want %v", err, ErrInvalidFarmHarvestInput)
			}
		})
	}
}

func TestFarmHarvestCreateScopesToOwnedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	service := NewFarmHarvestService(pool)
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	harvestID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_harvests.*WHERE cc.id = [$]1 AND fp.user_id = [$]2 RETURNING`).
		WithArgs(cycleID, userID, 125.5, "kg", "2026-10-19", "Grade A", "First pick").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "quantity", "unit", "harvest_date", "quality_grade", "notes"}).
			AddRow(harvestID, cycleID, 125.5, "kg", "2026-10-19", "Grade A", "First pick"))

	got, err := service.Create(context.Background(), userID, cycleID, CreateFarmHarvestInput{
		Quantity: 125.5, Unit: " kg ", HarvestDate: "2026-10-19",
		QualityGrade: " Grade A ", Notes: " First pick ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != harvestID || got.CropCycleID != cycleID || got.Quantity != 125.5 || got.Unit != "kg" {
		t.Fatalf("unexpected harvest: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFarmHarvestListHidesUnownedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	service := NewFarmHarvestService(pool)
	pool.ExpectQuery(`SELECT EXISTS.*WHERE cc.id = [$]1 AND fp.user_id = [$]2`).
		WithArgs("cycle-foreign", "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	_, err = service.List(context.Background(), "user-1", "cycle-foreign")
	if err != ErrFarmHarvestNotFound {
		t.Fatalf("List error = %v, want ErrFarmHarvestNotFound", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
