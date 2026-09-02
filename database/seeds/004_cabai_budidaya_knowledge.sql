-- M002.4: source-backed structural content seed.
-- This seed records a concise educational summary, not a replacement for the source document.
-- Source: Direktorat Sayuran dan Tanaman Obat, SOP Budidaya Cabai Rawit (2020).

INSERT INTO knowledge_documents (slug, title, summary, status, author_name)
VALUES (
  'cabai-rawit-budidaya-dasar',
  'Budidaya Cabai Rawit: Dasar Budidaya',
  'Panduan edukasi dasar budidaya cabai rawit dari persiapan benih dan lahan sampai pemeliharaan, pengendalian OPT, dan panen.',
  'draft',
  'Setyanto, Prihasto; Direktorat Sayuran dan Tanaman Obat'
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO knowledge_versions (document_id, version_no, content, source_id)
SELECT d.id, 1,
'Panduan edukasi dasar: budidaya cabai perlu dilakukan secara bertahap dan terencana, mulai dari penyediaan benih, persiapan lahan, penanaman, pengairan, pemeliharaan dan pemupukan, pemasangan penyangga tanaman, sanitasi, pengendalian organisme pengganggu tumbuhan (OPT), hingga panen. Pelaksanaan di lapangan perlu disesuaikan dengan kondisi lokasi dan pedoman teknis yang berlaku.',
s.id
FROM knowledge_documents d
JOIN knowledge_sources s ON s.title = 'Standar Operasional Prosedur Budidaya Cabai Rawit'
WHERE d.slug = 'cabai-rawit-budidaya-dasar'
  AND NOT EXISTS (SELECT 1 FROM knowledge_versions v WHERE v.document_id = d.id AND v.version_no = 1);

INSERT INTO knowledge_crop_links (document_id, crop_id)
SELECT d.id, c.id
FROM knowledge_documents d
JOIN crops c ON c.code = 'cabai'
WHERE d.slug = 'cabai-rawit-budidaya-dasar'
ON CONFLICT DO NOTHING;
