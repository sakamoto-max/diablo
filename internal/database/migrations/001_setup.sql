-- Write your migrate up statements here

CREATE OR REPLACE FUNCTION update_suite_last_changed()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    -- INSERT
    IF TG_OP = 'INSERT' THEN
        UPDATE suites
        SET last_changed = NOW()
        WHERE id = NEW.suite_id;

        RETURN NEW;
    END IF;

    -- UPDATE
    IF TG_OP = 'UPDATE' THEN
        -- Update the current suite
        UPDATE suites
        SET last_changed = NOW()
        WHERE id = NEW.suite_id;

        -- If suite_id itself changed, update the old suite too
        IF OLD.suite_id IS DISTINCT FROM NEW.suite_id THEN
            UPDATE suites
            SET last_changed = NOW()
            WHERE id = OLD.suite_id;
        END IF;

        RETURN NEW;
    END IF;

    -- DELETE
    IF TG_OP = 'DELETE' THEN
        UPDATE suites
        SET last_changed = NOW()
        WHERE id = OLD.suite_id;

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$;

---- create above / drop below ----

DROP FUNCTION IF EXISTS update_suite_last_changed;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
