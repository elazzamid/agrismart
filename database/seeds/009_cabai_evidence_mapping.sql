-- M002.7: conservative symptom-to-problem mappings.
-- These mappings are intended for candidate ranking only; they are not definitive diagnoses.

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Serangga kecil pada daun dapat menjadi petunjuk awal keberadaan hama pengisap; lakukan pemeriksaan lapangan sebelum menyimpulkan.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='serangga kecil pada daun' AND p.name='Thrips'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Daun mengeriting dapat berkaitan dengan aktivitas hama pengisap; periksa bagian bawah daun dan pucuk.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='daun mengeriting' AND p.name='Kutu daun persik'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, pest_id, weight, evidence_note)
SELECT s.id, p.id, 2, 'Serangga kecil pada daun dapat menjadi petunjuk keberadaan kutu kebul; lakukan pemeriksaan visual dan konfirmasi gejala terkait.'
FROM symptoms s CROSS JOIN pests p
WHERE s.name='serangga kecil pada daun' AND p.name='Kutu kebul'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 3, 'Bercak gelap pada buah dan pembusukan perlu diperiksa bersama bentuk dan perkembangan lesi untuk membedakan penyebab.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='bercak gelap pada buah' AND d.name='Antraknosa'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Buah membusuk merupakan gejala umum yang memerlukan pemeriksaan tambahan sebelum dikaitkan dengan antraknosa.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='buah membusuk' AND d.name='Antraknosa'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Daun menguning dapat muncul pada berbagai kondisi; untuk virus kuning perlu pemeriksaan gejala dan vektor terkait.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='daun menguning' AND d.name='Virus kuning'
ON CONFLICT DO NOTHING;

INSERT INTO problem_symptoms (symptom_id, disease_id, weight, evidence_note)
SELECT s.id, d.id, 2, 'Pertumbuhan terhambat dapat memiliki banyak penyebab; gunakan sebagai bukti pendukung, bukan diagnosis tunggal.'
FROM symptoms s CROSS JOIN diseases d
WHERE s.name='pertumbuhan terhambat' AND d.name='Virus kuning'
ON CONFLICT DO NOTHING;
