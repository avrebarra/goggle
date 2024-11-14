-- toggles definition
CREATE TABLE toggles (
    id TEXT(100) NOT NULL,
    status BOOLEAN,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_toggles_id ON toggles (id);

-- access_logs definition
CREATE TABLE access_logs (
    id INTEGER PRIMARY KEY,
    toggle_id TEXT,
    created_at DATETIME
);