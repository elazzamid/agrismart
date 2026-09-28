-- M002.4: source registry for the remaining MVP crops.
-- Registry only: no agronomic recommendations or pesticide instructions are copied here.
-- Sources are authoritative Indonesian agricultural institutional publications.

INSERT INTO knowledge_sources (title, publisher, source_url, source_type)
VALUES
('Budidaya Padi (Oryza sativa L.)', 'Balai Besar Pelatihan Pertanian Binuang', 'https://repository.pertanian.go.id/handle/123456789/27098', 'extension'),
('Prosedur Operasional Standar (POS) Budi Daya Padi Sawah', 'Pusat Penelitian dan Pengembangan Tanaman Pangan', 'https://repository.pertanian.go.id/handle/123456789/25727', 'government'),
('Budidaya Jagung Terstandar', 'Pertanian Press', 'https://repository.pertanian.go.id/handle/123456789/23686', 'government'),
('Budidaya Jagung', 'BPTP Sumatera Utara', 'https://repository.pertanian.go.id/handle/123456789/7188', 'extension')
ON CONFLICT DO NOTHING;
