-- Write your migrate up statements here

ALTER TABLE SUITES
ADD CONSTRAINT suites_name_unique UNIQUE (NAME);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
