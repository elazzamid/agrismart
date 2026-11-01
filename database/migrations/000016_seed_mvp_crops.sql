-- AgriSmart M004 demo prerequisite.
-- Seed only the MVP crop names explicitly defined in docs/MASTER_PLAN.md.
-- This catalog data does not imply agronomic recommendations or guidance.
INSERT INTO crops (code, name, description, is_active)
VALUES
    ('cabai', 'Cabai', NULL, TRUE),
    ('padi', 'Padi', NULL, TRUE),
    ('jagung', 'Jagung', NULL, TRUE)
ON CONFLICT (code) DO NOTHING;
