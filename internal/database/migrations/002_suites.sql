-- Write your migrate up statements here

CREATE TABLE SUITES(
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    NAME TEXT NOT NULL,
    LAST_CHANGED TIMESTAMPTZ NOT NULL
)


---- create above / drop below ----

DROP TABLE IF EXISTS SUITES;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

