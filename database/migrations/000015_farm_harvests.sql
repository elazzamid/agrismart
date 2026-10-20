-- AgriSmart M003.4
-- Records actual harvest outcomes; it does not predict yield.
CREATE TABLE farm_harvests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_cycle_id UUID NOT NULL REFERENCES crop_cycles(id) ON DELETE CASCADE,
    quantity NUMERIC(14,3) NOT NULL,
    unit TEXT NOT NULL,
    harvest_date DATE NOT NULL,
    quality_grade TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT farm_harvests_quantity_check CHECK (quantity >= 0.001 AND quantity <= 99999999999.999 AND quantity <> 'NaN'::numeric),
    CONSTRAINT farm_harvests_unit_check CHECK (length(btrim(unit)) BETWEEN 1 AND 20),
    CONSTRAINT farm_harvests_quality_grade_check CHECK (quality_grade IS NULL OR length(btrim(quality_grade)) BETWEEN 1 AND 40),
    CONSTRAINT farm_harvests_notes_check CHECK (notes IS NULL OR length(notes) <= 2000)
);

CREATE INDEX idx_farm_harvests_cycle_date ON farm_harvests(crop_cycle_id, harvest_date DESC, created_at DESC);
