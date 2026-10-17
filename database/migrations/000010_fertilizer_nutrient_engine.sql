-- AgriSmart M003.1
-- Nutrient requirements are guidance inputs; recommendations must remain source-backed.

CREATE TABLE nutrients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE crop_nutrient_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE CASCADE,
    growth_stage_id UUID REFERENCES crop_growth_stages(id) ON DELETE CASCADE,
    nutrient_id UUID NOT NULL REFERENCES nutrients(id) ON DELETE CASCADE,
    requirement_min NUMERIC(12,4),
    requirement_max NUMERIC(12,4),
    unit TEXT NOT NULL,
    source_document_id UUID REFERENCES knowledge_documents(id) ON DELETE SET NULL,
    source_version_id UUID REFERENCES knowledge_versions(id) ON DELETE SET NULL,
    notes TEXT,
    CHECK (requirement_min IS NULL OR requirement_min >= 0),
    CHECK (requirement_max IS NULL OR requirement_max >= 0),
    CHECK (requirement_min IS NULL OR requirement_max IS NULL OR requirement_min <= requirement_max),
    UNIQUE (crop_id, growth_stage_id, nutrient_id, source_document_id, source_version_id)
);

CREATE INDEX idx_crop_nutrient_requirements_crop ON crop_nutrient_requirements(crop_id);
CREATE INDEX idx_crop_nutrient_requirements_stage ON crop_nutrient_requirements(growth_stage_id);
CREATE INDEX idx_crop_nutrient_requirements_nutrient ON crop_nutrient_requirements(nutrient_id);
