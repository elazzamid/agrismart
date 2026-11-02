-- AgriSmart M004 demo prerequisite.
-- Seed structural variety and growth-stage catalog rows for the MVP crops.
-- These examples are catalog metadata, not agronomic recommendations.

INSERT INTO crop_varieties (crop_id, name, description)
SELECT c.id, v.name, v.description
FROM crops c
JOIN (VALUES
    ('padi', 'IR64', 'Contoh varietas padi.'),
    ('padi', 'Inpari 32', 'Contoh varietas padi.'),
    ('jagung', 'Bisi 18', 'Contoh varietas jagung.'),
    ('jagung', 'NK 212', 'Contoh varietas jagung.'),
    ('cabai', 'Cabai Merah Keriting', 'Contoh kelompok varietas cabai.'),
    ('cabai', 'Cabai Rawit', 'Contoh kelompok varietas cabai.')
) AS v(code, name, description) ON v.code = c.code
ON CONFLICT (crop_id, name) DO NOTHING;

INSERT INTO crop_growth_stages (crop_id, name, sequence_no, min_days, max_days, description)
SELECT c.id, v.name, v.sequence_no, v.min_days, v.max_days, v.description
FROM crops c
JOIN (VALUES
    ('padi', 'Persemaian', 1, 0, 21, 'Tahap awal pembibitan.'),
    ('padi', 'Vegetatif', 2, 22, 55, 'Pertumbuhan daun dan anakan.'),
    ('padi', 'Generatif', 3, 56, 100, 'Pembentukan malai hingga pengisian gabah.'),
    ('padi', 'Pematangan', 4, 101, 120, 'Pematangan gabah menuju panen.'),
    ('jagung', 'Perkecambahan', 1, 0, 7, 'Perkecambahan dan kemunculan tanaman.'),
    ('jagung', 'Vegetatif', 2, 8, 45, 'Pembentukan daun dan biomassa.'),
    ('jagung', 'Pembungaan', 3, 46, 70, 'Pembentukan bunga dan penyerbukan.'),
    ('jagung', 'Pengisian Biji', 4, 71, 110, 'Pengisian dan pematangan biji.'),
    ('cabai', 'Persemaian', 1, 0, 30, 'Pembibitan sebelum pindah tanam.'),
    ('cabai', 'Vegetatif', 2, 31, 60, 'Pembentukan batang, daun, dan cabang.'),
    ('cabai', 'Berbunga', 3, 61, 90, 'Pembentukan bunga dan awal pembentukan buah.'),
    ('cabai', 'Pembuahan', 4, 91, 140, 'Pembentukan dan pembesaran buah.'),
    ('cabai', 'Panen', 5, 141, 200, 'Periode produksi dan pemanenan.')
) AS v(code, name, sequence_no, min_days, max_days, description) ON v.code = c.code
ON CONFLICT (crop_id, sequence_no) DO NOTHING;
