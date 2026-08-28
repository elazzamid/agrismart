-- M002.4: source-backed educational knowledge.
-- No actionable pesticide dose, interval, tank mix, or brand recommendation is stored here.

INSERT INTO knowledge_documents (slug, title, summary, status, author_name)
VALUES
 ('cabai-opt-thrips','Cabai - Thrips dan Pemantauan','Mengenali gejala awal dan pentingnya pemantauan OPT pada tanaman cabai.','draft','Direktorat Sayuran dan Tanaman Obat / BPTP'),
 ('cabai-opt-antraknosa','Cabai - Antraknosa dan Pencegahan','Mengenali risiko antraknosa dan pendekatan pencegahan berbasis budidaya sehat.','draft','Direktorat Sayuran dan Tanaman Obat / BPTP'),
 ('cabai-opt-virus-kuning','Cabai - Virus Kuning dan Pencegahan','Mengenali gejala virus kuning dan pentingnya pengelolaan tanaman serta vektor.','draft','Direktorat Sayuran dan Tanaman Obat / BPTP')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO knowledge_versions (document_id, version_no, content, source_id)
SELECT d.id, 1,
CASE d.slug
 WHEN 'cabai-opt-thrips' THEN 'Thrips termasuk OPT penting pada cabai. Pemantauan tanaman secara rutin membantu menemukan gejala dan perubahan populasi lebih awal. Pengelolaan perlu mengikuti prinsip PHT dan pedoman resmi yang sesuai kondisi lapangan.'
 WHEN 'cabai-opt-antraknosa' THEN 'Antraknosa merupakan penyakit penting pada cabai dan dapat menyebabkan kerusakan pada buah. Pencegahan perlu menekankan sanitasi, bahan tanam sehat, pengelolaan lingkungan dan pemantauan rutin. Tindakan pengendalian harus mengikuti sumber resmi yang sesuai kondisi setempat.'
 WHEN 'cabai-opt-virus-kuning' THEN 'Virus kuning dapat menurunkan pertumbuhan dan produktivitas cabai. Pengelolaan perlu menekankan pencegahan, pengelolaan tanaman dan pengendalian vektor berdasarkan panduan resmi serta hasil pengamatan lapangan.'
END,
s.id
FROM knowledge_documents d
JOIN knowledge_sources s ON s.title = 'Petunjuk Teknis Pengendalian Hama dan Penyakit Utama Tanaman Cabai Merah dan Tomat'
WHERE d.slug IN ('cabai-opt-thrips','cabai-opt-antraknosa','cabai-opt-virus-kuning')
AND NOT EXISTS (SELECT 1 FROM knowledge_versions v WHERE v.document_id=d.id AND v.version_no=1);

INSERT INTO knowledge_pest_links (document_id, pest_id)
SELECT d.id, p.id FROM knowledge_documents d JOIN pests p ON p.name='Thrips'
WHERE d.slug='cabai-opt-thrips' ON CONFLICT DO NOTHING;

INSERT INTO knowledge_disease_links (document_id, disease_id)
SELECT d.id, p.id FROM knowledge_documents d JOIN diseases p ON p.name='Antraknosa'
WHERE d.slug='cabai-opt-antraknosa' ON CONFLICT DO NOTHING;

INSERT INTO knowledge_disease_links (document_id, disease_id)
SELECT d.id, p.id FROM knowledge_documents d JOIN diseases p ON p.name='Virus kuning'
WHERE d.slug='cabai-opt-virus-kuning' ON CONFLICT DO NOTHING;
