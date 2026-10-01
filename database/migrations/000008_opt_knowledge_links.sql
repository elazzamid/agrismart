CREATE TABLE knowledge_pest_links (
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    pest_id UUID NOT NULL REFERENCES pests(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, pest_id)
);

CREATE TABLE knowledge_disease_links (
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    disease_id UUID NOT NULL REFERENCES diseases(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, disease_id)
);

CREATE INDEX idx_knowledge_pest_links_pest_id ON knowledge_pest_links(pest_id);
CREATE INDEX idx_knowledge_disease_links_disease_id ON knowledge_disease_links(disease_id);
