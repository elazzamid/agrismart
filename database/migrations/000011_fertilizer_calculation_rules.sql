-- M003.3: keep fertilizer composition data mathematically safe.
ALTER TABLE fertilizer_nutrients
    ADD CONSTRAINT fertilizer_nutrients_percentage_max_check CHECK (percentage IS NULL OR percentage <= 100);

CREATE UNIQUE INDEX idx_fertilizer_nutrients_unique_code
    ON fertilizer_nutrients(fertilizer_id, nutrient_code);

CREATE INDEX idx_fertilizer_nutrients_code
    ON fertilizer_nutrients(nutrient_code);
