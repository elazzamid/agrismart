-- AgriSmart M001.2

CREATE TABLE crops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE crop_varieties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT crop_varieties_crop_name_unique UNIQUE (crop_id, name)
);

CREATE TABLE crop_growth_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sequence_no INTEGER NOT NULL,
    min_days INTEGER,
    max_days INTEGER,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT crop_growth_stages_sequence_unique UNIQUE (crop_id, sequence_no),
    CONSTRAINT crop_growth_stages_days_check CHECK (
        (min_days IS NULL OR min_days >= 0) AND
        (max_days IS NULL OR max_days >= 0) AND
        (min_days IS NULL OR max_days IS NULL OR min_days <= max_days)
    )
);

CREATE TABLE crop_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plot_id UUID NOT NULL REFERENCES farm_plots(id) ON DELETE CASCADE,
    crop_id UUID NOT NULL REFERENCES crops(id),
    variety_id UUID REFERENCES crop_varieties(id),
    planting_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT crop_cycles_status_check CHECK (status IN ('planned', 'active', 'harvested', 'cancelled'))
);

CREATE INDEX idx_crop_varieties_crop_id ON crop_varieties(crop_id);
CREATE INDEX idx_crop_growth_stages_crop_id ON crop_growth_stages(crop_id);
CREATE INDEX idx_crop_cycles_plot_id ON crop_cycles(plot_id);
CREATE INDEX idx_crop_cycles_crop_id ON crop_cycles(crop_id);
CREATE INDEX idx_crop_cycles_status ON crop_cycles(status);
