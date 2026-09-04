CREATE TABLE suites(
    ID INTEGER PRIMARY KEY,
    NAME TEXT NOT NULL UNIQUE,
    LAST_CHANGED TIMESTAMPTZ NOT NULL
);

--------------------------------------------


CREATE TABLE FILES(
    ID INTEGER PRIMARY KEY,
    SUITE_ID INTEGER NOT NULL REFERENCES suites(ID),
    PATH TEXT NOT NULL,
    DATA BLOB,
    FILE_TYPE TEXT NOT NULL
);

CREATE TRIGGER update_suite_last_changed_insert
AFTER INSERT ON FILES
FOR EACH ROW
BEGIN
    UPDATE suites
    SET last_changed = CURRENT_TIMESTAMP
    WHERE id = NEW.suite_id;
END;

CREATE TRIGGER update_suite_last_changed_update
AFTER UPDATE ON FILES
FOR EACH ROW
BEGIN
    UPDATE suites
    SET last_changed = CURRENT_TIMESTAMP
    WHERE id = NEW.suite_id;
END;

CREATE TRIGGER update_suite_last_changed_delete
AFTER DELETE ON FILES
FOR EACH ROW
BEGIN
    UPDATE suites
    SET last_changed = CURRENT_TIMESTAMP
    WHERE id = OLD.suite_id;
END;

--------------------------------------------

CREATE TABLE USER_IPS(
    ID INTEGER PRIMARY KEY,
    SUITE_ID INTEGER NOT NULL REFERENCES suites(ID),
    IP TEXT NOT NULL,
    LAST_SYNCED TIMESTAMPTZ NOT NULL
);

--------------------------------------------

CREATE TABLE EVENTS(
    ID INTEGER PRIMARY KEY,
    SUITE_ID INTEGER NOT NULL REFERENCES suites(ID),
    NAME TEXT NOT NULL,
    PATH TEXT NOT NULL,
    CHANGED_AT TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    RENAMED_TO TEXT
)
--------------------------------------------








