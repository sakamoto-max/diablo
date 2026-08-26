-- Write your migrate up statements here

CREATE TABLE FILES(
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    SUITE_ID UUID NOT NULL,
    PATH TEXT NOT NULL,
    DATA BYTEA,
    FILE_TYPE TEXT NOT NULL,
    FOREIGN KEY (SUITE_ID) REFERENCES SUITES(ID)
);

CREATE TRIGGER trg_update_suite_last_changed
AFTER INSERT OR UPDATE OR DELETE
ON files
FOR EACH ROW
EXECUTE FUNCTION update_suite_last_changed();


---- create above / drop below ----

DROP TABLE IF EXISTS FILES;

DROP TRIGGER IF EXISTS trg_update_suite_last_changed ON files;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
