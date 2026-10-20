package farm

import (
	"math"
	"testing"
)

func TestNormalizeFarmInput(t *testing.T) {
	valid := CreateFarmInputInput{
		InputType: " Fertilizer ",
		ProductName: " Urea ",
		Quantity: 1.25,
		Unit: " kg ",
		AppliedDate: "2026-10-19",
		Notes: " applied ",
	}
	got, err := normalizeFarmInput(valid)
	if err != nil {
		t.Fatalf("normalizeFarmInput(valid) error = %v", err)
	}
	if got.InputType != "fertilizer" || got.ProductName != "Urea" || got.Unit != "kg" || got.Notes != "applied" {
		t.Fatalf("unexpected normalized input: %#v", got)
	}
}

func TestNormalizeFarmInputRejectsInvalidValues(t *testing.T) {
	base := CreateFarmInputInput{
		InputType: "fertilizer",
		ProductName: "Urea",
		Quantity: 1,
		Unit: "kg",
		AppliedDate: "2026-10-19",
	}
	tests := []struct {
		name string
		mutate func(*CreateFarmInputInput)
	}{
		{name: "unknown type", mutate: func(in *CreateFarmInputInput) { in.InputType = "medicine" }},
		{name: "empty product", mutate: func(in *CreateFarmInputInput) { in.ProductName = " " }},
		{name: "zero quantity", mutate: func(in *CreateFarmInputInput) { in.Quantity = 0 }},
		{name: "negative quantity", mutate: func(in *CreateFarmInputInput) { in.Quantity = -1 }},
		{name: "NaN quantity", mutate: func(in *CreateFarmInputInput) { in.Quantity = math.NaN() }},
		{name: "infinite quantity", mutate: func(in *CreateFarmInputInput) { in.Quantity = math.Inf(1) }},
		{name: "empty unit", mutate: func(in *CreateFarmInputInput) { in.Unit = "" }},
		{name: "invalid date", mutate: func(in *CreateFarmInputInput) { in.AppliedDate = "2026-02-30" }},
		{name: "noncanonical date", mutate: func(in *CreateFarmInputInput) { in.AppliedDate = "2026-1-9" }},
		{name: "notes too long", mutate: func(in *CreateFarmInputInput) { in.Notes = string(make([]byte, 2001)) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := base
			tt.mutate(&in)
			if _, err := normalizeFarmInput(in); err != ErrInvalidFarmInput {
				t.Fatalf("normalizeFarmInput() error = %v, want %v", err, ErrInvalidFarmInput)
			}
		})
	}
}
