-- AgriSmart M003.3
-- Records farm operating costs; this is accounting history, not financial advice.
CREATE TABLE farm_expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_cycle_id UUID NOT NULL REFERENCES crop_cycles(id) ON DELETE CASCADE,
    expense_type TEXT NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    expense_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT farm_expenses_type_check CHECK (expense_type IN ('seed', 'fertilizer', 'pesticide', 'labor', 'irrigation', 'equipment', 'transport', 'other')),
    CONSTRAINT farm_expenses_description_check CHECK (length(btrim(description)) BETWEEN 1 AND 160),
    CONSTRAINT farm_expenses_amount_check CHECK (amount >= 0.01 AND amount <= 999999999999.99 AND amount <> 'NaN'::numeric),
    CONSTRAINT farm_expenses_notes_check CHECK (notes IS NULL OR length(notes) <= 2000)
);

CREATE INDEX idx_farm_expenses_cycle_date ON farm_expenses(crop_cycle_id, expense_date DESC, created_at DESC);
