-- AgriSmart M003.2
-- Records what a farmer actually used; this table is not an agronomic recommendation.
CREATE TABLE farm_inputs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_cycle_id UUID NOT NULL REFERENCES crop_cycles(id) ON DELETE CASCADE,
    input_type TEXT NOT NULL,
    product_name TEXT NOT NULL,
    quantity NUMERIC(14,3) NOT NULL,
    unit TEXT NOT NULL,
    applied_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT farm_inputs_type_check CHECK (input_type IN ('fertilizer', 'pesticide', 'seed', 'other')),
    CONSTRAINT farm_inputs_product_name_check CHECK (length(btrim(product_name)) BETWEEN 1 AND 120),
    CONSTRAINT farm_inputs_quantity_check CHECK (quantity >= 0.001),
    CONSTRAINT farm_inputs_unit_check CHECK (length(btrim(unit)) BETWEEN 1 AND 20),
    CONSTRAINT farm_inputs_notes_check CHECK (notes IS NULL OR length(notes) <= 2000)
);

CREATE INDEX idx_farm_inputs_cycle_date ON farm_inputs(crop_cycle_id, applied_date DESC, created_at DESC);
