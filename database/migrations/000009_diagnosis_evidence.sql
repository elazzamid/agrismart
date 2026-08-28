CREATE TABLE symptoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE problem_symptoms (
    symptom_id UUID NOT NULL REFERENCES symptoms(id) ON DELETE CASCADE,
    pest_id UUID REFERENCES pests(id) ON DELETE CASCADE,
    disease_id UUID REFERENCES diseases(id) ON DELETE CASCADE,
    weight INTEGER NOT NULL DEFAULT 1,
    evidence_note TEXT,
    PRIMARY KEY (symptom_id, pest_id, disease_id),
    CONSTRAINT problem_symptoms_one_problem CHECK ((pest_id IS NOT NULL) <> (disease_id IS NOT NULL)),
    CONSTRAINT problem_symptoms_weight_check CHECK (weight > 0)
);

CREATE INDEX idx_problem_symptoms_pest_id ON problem_symptoms(pest_id);
CREATE INDEX idx_problem_symptoms_disease_id ON problem_symptoms(disease_id);
