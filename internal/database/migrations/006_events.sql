-- Write your migrate up statements here


CREATE TABLE EVENTS(
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    SUITE_ID UUID NOT NULL,
    NAME TEXT NOT NULL,
    PATH TEXT NOT NULL,
    CHANGED_AT TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    RENAMED_TO TEXT,
    FOREIGN KEY (SUITE_ID) REFERENCES SUITES(ID)
)

---- create above / drop below ----

DROP TABLE IF EXISTS EVENTS;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
