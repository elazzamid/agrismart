-- M002.4: source-backed educational IPM content.
-- Source: Direktorat Jenderal Hortikultura / BPTP Jambi, Hama dan Penyakit pada Tanaman Cabai serta Pengendaliannya.
-- Actionable pesticide details are intentionally not embedded here.

INSERT INTO knowledge_documents (slug, title, summary, status, author_name)
VALUES (
  'cabai-pht-dasar',
  'Cabai: Pengendalian OPT Terpadu',
  'Edukasi dasar pengamatan dan pengendalian organisme pengganggu tanaman cabai dengan pendekatan Pengendalian Hama Terpadu (PHT).',
  'draft',
  'Direktorat Jenderal Hortikultura / BPTP Jambi'
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO knowledge_versions (document_id, version_no, content, source_id)
SELECT d.id, 1,
'Pengendalian OPT pada cabai perlu diawali pengamatan berkala untuk mengenali jenis OPT serta luas dan intensitas serangan. Tindakan pengendalian dipilih berdasarkan kondisi pertanaman dan prinsip Pengendalian Hama Terpadu (PHT). Cara non-kimia dan tindakan budidaya perlu dipertimbangkan terlebih dahulu; pestisida bukan pilihan pertama dan penggunaannya harus mengikuti ketentuan label serta aturan yang berlaku.',
s.id
FROM knowledge_documents d
JOIN knowledge_sources s ON s.title = 'HAMA DAN PENYAKIT PADA TANAMAN CABAI SERTA PENGENDALIANNYA'
WHERE d.slug = 'cabai-pht-dasar'
  AND NOT EXISTS (SELECT 1 FROM knowledge_versions v WHERE v.document_id = d.id AND v.version_no = 1);

INSERT INTO knowledge_crop_links (document_id, crop_id)
SELECT d.id, c.id
FROM knowledge_documents d
JOIN crops c ON c.code = 'cabai'
WHERE d.slug = 'cabai-pht-dasar'
ON CONFLICT DO NOTHING;
