-- Educational knowledge only. Actionable pesticide recommendations are intentionally excluded.
INSERT INTO knowledge_documents (slug, title, summary, status, author_name)
VALUES
 ('cabai-opt-thrips','Cabai - Thrips dan Pemantauan','Mengenali gejala awal dan pentingnya pemantauan OPT pada tanaman cabai.','draft','AgriSmart editorial'),
 ('cabai-opt-antraknosa','Cabai - Antraknosa dan Pencegahan','Mengenali risiko antraknosa dan pendekatan pencegahan berbasis budidaya sehat.','draft','AgriSmart editorial'),
 ('cabai-opt-virus-kuning','Cabai - Virus Kuning dan Pencegahan','Mengenali gejala virus kuning dan pentingnya pengelolaan tanaman serta vektor.','draft','AgriSmart editorial')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO knowledge_versions (document_id, version_no, content, source_id)
SELECT d.id, 1,
CASE d.slug
 WHEN 'cabai-opt-thrips' THEN 'Thrips termasuk OPT penting pada cabai. Pemantauan tanaman secara rutin membantu menemukan gejala dan perubahan populasi lebih awal. Pengelolaan sebaiknya mengikuti prinsip PHT: tanaman sehat, pemantauan rutin, pemanfaatan pengendalian non-kimia dan musuh alami bila tersedia, serta penggunaan pestisida hanya ketika diperlukan dan sesuai sumber resmi.'
 WHEN 'cabai-opt-antraknosa' THEN 'Antraknosa merupakan penyakit penting pada cabai dan dapat menyebabkan kerusakan pada buah. Pencegahan perlu menekankan sanitasi, penggunaan bahan tanam sehat, pengelolaan lingkungan dan pemantauan rutin. Tindakan pengendalian harus mengikuti rekomendasi sumber resmi yang sesuai dengan kondisi setempat.'
 WHEN 'cabai-opt-virus-kuning' THEN 'Virus kuning dapat menurunkan pertumbuhan dan produktivitas cabai. Karena penyakit virus tidak ditangani dengan cara yang sama seperti penyakit jamur, pengelolaan perlu menekankan pencegahan, pengelolaan tanaman dan pengendalian vektor berdasarkan panduan resmi serta hasil pengamatan lapangan.'
END,
s.id
FROM knowledge_documents d
JOIN knowledge_sources s ON s.title = 'Petunjuk Teknis Pengendalian Hama dan Penyakit Utama Tanaman Cabai Merah dan Tomat'
WHERE d.slug IN ('cabai-opt-thrips','cabai-opt-antraknosa','cabai-opt-virus-kuning')
AND NOT EXISTS (SELECT 1 FROM knowledge_versions v WHERE v.document_id=d.id AND v.version_no=1);
