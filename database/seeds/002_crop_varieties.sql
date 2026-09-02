-- Structural variety examples only. They are not recommendations.
-- The catalog can be expanded after source validation.

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
