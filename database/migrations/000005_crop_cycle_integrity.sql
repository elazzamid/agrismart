-- AgriSmart M001.6 integrity constraint.
-- A crop cycle's optional variety must belong to the same crop.

CREATE OR REPLACE FUNCTION validate_crop_cycle_variety()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.variety_id IS NOT NULL AND NOT EXISTS (
        SELECT 1 FROM crop_varieties cv
        WHERE cv.id = NEW.variety_id AND cv.crop_id = NEW.crop_id
    ) THEN
        RAISE EXCEPTION 'crop variety does not belong to selected crop';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_crop_cycle_variety
BEFORE INSERT OR UPDATE OF crop_id, variety_id ON crop_cycles
FOR EACH ROW EXECUTE FUNCTION validate_crop_cycle_variety();
