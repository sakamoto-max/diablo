-- Write your migrate up statements here


CREATE TABLE USER_IPS(
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    SUITE_ID UUID NOT NULL,
    IP TEXT NOT NULL,
    LAST_SYNCED TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (SUITE_ID) REFERENCES SUITES(ID)
)

---- create above / drop below ----

DROP TABLE IF EXISTS USER_IPS;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
