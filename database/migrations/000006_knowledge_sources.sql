-- AgriSmart M002.1
-- Knowledge must be traceable and reviewable before publication.

CREATE TABLE knowledge_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    publisher TEXT,
    source_url TEXT,
    source_type TEXT NOT NULL DEFAULT 'reference',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT knowledge_sources_type_check CHECK (source_type IN ('government', 'research', 'extension', 'label', 'reference'))
);

CREATE TABLE knowledge_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    summary TEXT,
    status TEXT NOT NULL DEFAULT 'draft',
    author_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT knowledge_documents_status_check CHECK (status IN ('draft', 'review', 'validated', 'published', 'archived'))
);

CREATE TABLE knowledge_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES knowledge_documents(id) ON DELETE CASCADE,
    version_no INTEGER NOT NULL,
    content TEXT NOT NULL,
    source_id UUID REFERENCES knowledge_sources(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT knowledge_versions_unique UNIQUE (document_id, version_no),
    CONSTRAINT knowledge_versions_version_check CHECK (version_no > 0)
);

CREATE TABLE knowledge_validations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id UUID NOT NULL REFERENCES knowledge_versions(id) ON DELETE CASCADE,
    validator_name TEXT NOT NULL,
    decision TEXT NOT NULL,
    notes TEXT,
    validated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT knowledge_validations_decision_check CHECK (decision IN ('approved', 'rejected', 'needs_revision'))
);

CREATE INDEX idx_knowledge_versions_document_id ON knowledge_versions(document_id);
CREATE INDEX idx_knowledge_validations_version_id ON knowledge_validations(version_id);
