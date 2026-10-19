-- Source registry only. Content is not copied into seed data.
-- The source records point to authoritative Indonesian agricultural references.

INSERT INTO knowledge_sources (title, publisher, source_url, source_type)
VALUES
('Budidaya Tanaman Cabai', 'BPTP Sumatera Utara', 'https://repository.pertanian.go.id/handle/123456789/7189', 'extension'),
('Standar Operasional Prosedur Budidaya Cabai Rawit', 'Direktorat Sayuran dan Tanaman Obat', 'https://repository.pertanian.go.id/handle/123456789/11340', 'government'),
('Hama dan Penyakit pada Tanaman Cabai serta Pengendaliannya', 'Direktorat Jenderal Hortikultura', 'https://hortikultura.pertanian.go.id/hama-dan-penyakit-pada-tanaman-cabai-serta-pengendaliannya/', 'government'),
('Petunjuk Teknis Pengendalian Hama dan Penyakit Utama Tanaman Cabai Merah dan Tomat', 'BPTP Jawa Barat', 'https://repository.pertanian.go.id/handle/123456789/14046', 'extension')
ON CONFLICT DO NOTHING;
