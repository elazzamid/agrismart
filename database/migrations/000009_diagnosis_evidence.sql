CREATE TABLE symptoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- A symptom may map to either a pest or a disease. Keep the nullable
-- foreign-key columns out of a primary key: PostgreSQL primary-key columns
-- are implicitly NOT NULL, which would conflict with the one-problem CHECK.
CREATE TABLE problem_symptoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symptom_id UUID NOT NULL REFERENCES symptoms(id) ON DELETE CASCADE,
    pest_id UUID REFERENCES pests(id) ON DELETE CASCADE,
    disease_id UUID REFERENCES diseases(id) ON DELETE CASCADE,
    weight INTEGER NOT NULL DEFAULT 1,
    evidence_note TEXT,
    CONSTRAINT problem_symptoms_one_problem CHECK ((pest_id IS NOT NULL) <> (disease_id IS NOT NULL)),
    CONSTRAINT problem_symptoms_weight_check CHECK (weight > 0)
);

CREATE UNIQUE INDEX idx_problem_symptoms_unique_pest
    ON problem_symptoms(symptom_id, pest_id)
    WHERE pest_id IS NOT NULL;
CREATE UNIQUE INDEX idx_problem_symptoms_unique_disease
    ON problem_symptoms(symptom_id, disease_id)
    WHERE disease_id IS NOT NULL;

-- Each evidence row is tied to a specific source document version. As above,
-- nullable pest/disease IDs must not be primary-key columns.
CREATE TABLE problem_symptom_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symptom_id UUID NOT NULL REFERENCES symptoms(id) ON DELETE CASCADE,
    pest_id UUID REFERENCES pests(id) ON DELETE CASCADE,
    disease_id UUID REFERENCES diseases(id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    version_id UUID NOT NULL REFERENCES knowledge_versions(id) ON DELETE CASCADE,
    CONSTRAINT problem_symptom_sources_one_problem CHECK ((pest_id IS NOT NULL) <> (disease_id IS NOT NULL))
);

CREATE UNIQUE INDEX idx_problem_symptom_sources_unique_pest
    ON problem_symptom_sources(symptom_id, pest_id, document_id, version_id)
    WHERE pest_id IS NOT NULL;
CREATE UNIQUE INDEX idx_problem_symptom_sources_unique_disease
    ON problem_symptom_sources(symptom_id, disease_id, document_id, version_id)
    WHERE disease_id IS NOT NULL;

CREATE INDEX idx_problem_symptoms_pest_id ON problem_symptoms(pest_id);
CREATE INDEX idx_problem_symptoms_disease_id ON problem_symptoms(disease_id);
CREATE INDEX idx_problem_symptom_sources_document_id ON problem_symptom_sources(document_id);
CREATE INDEX idx_problem_symptom_sources_version_id ON problem_symptom_sources(version_id);
