-- AgriSmart M002.1

CREATE TABLE cultivation_guides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE CASCADE,
    growth_stage_id UUID REFERENCES crop_growth_stages(id) ON DELETE SET NULL,
    document_id UUID NOT NULL REFERENCES knowledge_documents(id),
    title TEXT NOT NULL,
    instructions TEXT NOT NULL
);

CREATE TABLE fertilizers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    formulation TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fertilizer_nutrients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fertilizer_id UUID NOT NULL REFERENCES fertilizers(id) ON DELETE CASCADE,
    nutrient_code TEXT NOT NULL,
    percentage NUMERIC(6,3),
    CONSTRAINT fertilizer_nutrients_percentage_check CHECK (percentage IS NULL OR percentage >= 0)
);

CREATE TABLE active_ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE pesticides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    formulation TEXT,
    active_ingredient_id UUID REFERENCES active_ingredients(id),
    registration_number TEXT,
    label_source_id UUID REFERENCES knowledge_sources(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    scientific_name TEXT,
    description TEXT
);

CREATE TABLE diseases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    scientific_name TEXT,
    description TEXT
);

CREATE TABLE weeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    scientific_name TEXT,
    description TEXT
);

CREATE TABLE knowledge_crop_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    crop_id UUID NOT NULL REFERENCES crops(id) ON DELETE CASCADE,
    variety_id UUID REFERENCES crop_varieties(id) ON DELETE SET NULL,
    growth_stage_id UUID REFERENCES crop_growth_stages(id) ON DELETE SET NULL,
    UNIQUE (document_id, crop_id, variety_id, growth_stage_id)
);

CREATE INDEX idx_cultivation_guides_crop_id ON cultivation_guides(crop_id);
CREATE INDEX idx_cultivation_guides_growth_stage_id ON cultivation_guides(growth_stage_id);
CREATE INDEX idx_pesticides_active_ingredient_id ON pesticides(active_ingredient_id);
CREATE INDEX idx_knowledge_crop_links_crop_id ON knowledge_crop_links(crop_id);
