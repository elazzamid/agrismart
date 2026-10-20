-- AgriSmart M003.1
CREATE TABLE farm_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_cycle_id UUID NOT NULL REFERENCES crop_cycles(id) ON DELETE CASCADE,
    activity_type TEXT NOT NULL,
    activity_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT farm_activities_type_check CHECK (length(btrim(activity_type)) BETWEEN 1 AND 80),
    CONSTRAINT farm_activities_notes_check CHECK (notes IS NULL OR length(notes) <= 2000)
);
CREATE INDEX idx_farm_activities_cycle_date ON farm_activities(crop_cycle_id, activity_date DESC, created_at DESC);
