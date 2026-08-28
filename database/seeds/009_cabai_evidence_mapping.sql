-- M002.7: conservative symptom-to-problem mappings with provenance.
-- Mappings are candidate-ranking signals, not definitive diagnoses.

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Serangga kecil pada daun dapat menjadi petunjuk awal keberadaan hama pengisap; lakukan pemeriksaan lapangan sebelum menyimpulkan.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='serangga kecil pada daun' AND p.name='Thrips' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Daun mengeriting dapat berkaitan dengan aktivitas hama pengisap; periksa bagian bawah daun dan pucuk.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='daun mengeriting' AND p.name='Kutu daun persik' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Serangga kecil pada daun dapat menjadi petunjuk keberadaan kutu kebul; lakukan pemeriksaan visual dan konfirmasi gejala terkait.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='serangga kecil pada daun' AND p.name='Kutu kebul' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 3, 'Bercak gelap pada buah dan pembusukan perlu diperiksa bersama bentuk dan perkembangan lesi untuk membedakan penyebab.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='bercak gelap pada buah' AND d.name='Antraknosa' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Buah membusuk merupakan gejala umum yang memerlukan pemeriksaan tambahan sebelum dikaitkan dengan antraknosa.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='buah membusuk' AND d.name='Antraknosa' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Daun menguning dapat muncul pada berbagai kondisi; untuk virus kuning perlu pemeriksaan gejala dan vektor terkait.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='daun menguning' AND d.name='Virus kuning' ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Pertumbuhan terhambat dapat memiliki banyak penyebab; gunakan sebagai bukti pendukung, bukan diagnosis tunggal.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='pertumbuhan terhambat' AND d.name='Virus kuning' ON CONFLICT DO NOTHING;

-- Attach each mapping to the exact approved knowledge version that supports it.
INSERT INTO problem_symptom_sources (symptom_id, pest_id, document_id, version_id)
SELECT s.id, p.id, d.id, v.id
FROM symptoms s JOIN pests p ON p.name='Thrips'
JOIN knowledge_documents d ON d.slug='cabai-opt-thrips'
JOIN knowledge_versions v ON v.document_id=d.id AND v.version_no=1
WHERE s.name='serangga kecil pada daun'
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptom_sources (symptom_id, pest_id, document_id, version_id)
SELECT s.id, p.id, d.id, v.id
FROM symptoms s JOIN pests p ON p.name='Kutu daun persik'
JOIN knowledge_documents d ON d.slug='cabai-opt-thrips'
JOIN knowledge_versions v ON v.document_id=d.id AND v.version_no=1
WHERE s.name='daun mengeriting'
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptom_sources (symptom_id, pest_id, document_id, version_id)
SELECT s.id, p.id, d.id, v.id
FROM symptoms s JOIN pests p ON p.name='Kutu kebul'
JOIN knowledge_documents d ON d.slug='cabai-opt-virus-kuning'
JOIN knowledge_versions v ON v.document_id=d.id AND v.version_no=1
WHERE s.name='serangga kecil pada daun'
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptom_sources (symptom_id, disease_id, document_id, version_id)
SELECT s.id, dse.id, d.id, v.id
FROM symptoms s JOIN diseases dse ON dse.name='Antraknosa'
JOIN knowledge_documents d ON d.slug='cabai-opt-antraknosa'
JOIN knowledge_versions v ON v.document_id=d.id AND v.version_no=1
WHERE s.name IN ('bercak gelap pada buah','buah membusuk')
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptom_sources (symptom_id, disease_id, document_id, version_id)
SELECT s.id, dse.id, d.id, v.id
FROM symptoms s JOIN diseases dse ON dse.name='Virus kuning'
JOIN knowledge_documents d ON d.slug='cabai-opt-virus-kuning'
JOIN knowledge_versions v ON v.document_id=d.id AND v.version_no=1
WHERE s.name IN ('daun menguning','pertumbuhan terhambat')
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
ON CONFLICT DO NOTHING;
