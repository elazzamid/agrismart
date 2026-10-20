package farm

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestNormalizeFarmExpenseInput(t *testing.T) {
	got, err := normalizeFarmExpenseInput(CreateFarmExpenseInput{
		ExpenseType: " Labor ",
		Description: " Harvest help ",
		Amount: 250000,
		ExpenseDate: "2026-10-19",
		Notes: " Paid in cash ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ExpenseType != "labor" || got.Description != "Harvest help" || got.Notes != "Paid in cash" {
		t.Fatalf("unexpected normalized expense: %+v", got)
	}
}

func TestNormalizeFarmExpenseInputRejectsInvalidValues(t *testing.T) {
	base := CreateFarmExpenseInput{
		ExpenseType: "labor", Description: "Harvest help", Amount: 100,
		ExpenseDate: "2026-10-19",
	}
	tests := []struct {
		name string
		mutate func(*CreateFarmExpenseInput)
	}{
		{name: "unknown type", mutate: func(in *CreateFarmExpenseInput) { in.ExpenseType = "loan" }},
		{name: "empty description", mutate: func(in *CreateFarmExpenseInput) { in.Description = " " }},
		{name: "zero amount", mutate: func(in *CreateFarmExpenseInput) { in.Amount = 0 }},
		{name: "negative amount", mutate: func(in *CreateFarmExpenseInput) { in.Amount = -1 }},
		{name: "NaN amount", mutate: func(in *CreateFarmExpenseInput) { in.Amount = math.NaN() }},
		{name: "infinite amount", mutate: func(in *CreateFarmExpenseInput) { in.Amount = math.Inf(1) }},
		{name: "amount too large", mutate: func(in *CreateFarmExpenseInput) { in.Amount = 1000000000000 }},
		{name: "invalid date", mutate: func(in *CreateFarmExpenseInput) { in.ExpenseDate = "2026-02-30" }},
		{name: "noncanonical date", mutate: func(in *CreateFarmExpenseInput) { in.ExpenseDate = "2026-1-9" }},
		{name: "notes too long", mutate: func(in *CreateFarmExpenseInput) { in.Notes = strings.Repeat("x", 2001) }},
		{name: "description too long", mutate: func(in *CreateFarmExpenseInput) { in.Description = strings.Repeat("x", 161) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := base
			tt.mutate(&in)
			if _, err := normalizeFarmExpenseInput(in); err != ErrInvalidFarmExpenseInput {
				t.Fatalf("normalizeFarmExpenseInput() error = %v, want %v", err, ErrInvalidFarmExpenseInput)
			}
		})
	}
}

func TestFarmExpenseCreateScopesToOwnedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	service := NewFarmExpenseService(pool)
	userID := "00000000-0000-0000-0000-000000000001"
	cycleID := "00000000-0000-0000-0000-000000000002"
	expenseID := "00000000-0000-0000-0000-000000000003"
	pool.ExpectQuery(`INSERT INTO farm_expenses.*WHERE cc.id = [$]1 AND fp.user_id = [$]2 RETURNING`).
		WithArgs(cycleID, userID, "labor", "Harvest help", 250000.0, "2026-10-19", "Paid in cash").
		WillReturnRows(pgxmock.NewRows([]string{"id", "crop_cycle_id", "expense_type", "description", "amount", "expense_date", "notes"}).
			AddRow(expenseID, cycleID, "labor", "Harvest help", 250000.0, "2026-10-19", "Paid in cash"))

	got, err := service.Create(context.Background(), userID, cycleID, CreateFarmExpenseInput{
		ExpenseType: " Labor ", Description: " Harvest help ", Amount: 250000,
		ExpenseDate: "2026-10-19", Notes: " Paid in cash ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != expenseID || got.CropCycleID != cycleID || got.ExpenseType != "labor" || got.Amount != 250000 {
		t.Fatalf("unexpected expense: %+v", got)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFarmExpenseListHidesUnownedCropCycle(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	service := NewFarmExpenseService(pool)
	pool.ExpectQuery(`SELECT EXISTS.*WHERE cc.id = [$]1 AND fp.user_id = [$]2`).
		WithArgs("cycle-foreign", "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	_, err = service.List(context.Background(), "user-1", "cycle-foreign")
	if err != ErrFarmExpenseNotFound {
		t.Fatalf("List error = %v, want ErrFarmExpenseNotFound", err)
	}
	if err := pool.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
