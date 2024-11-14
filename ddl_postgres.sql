-- toggles definition
CREATE TABLE toggles (
    id TEXT NOT NULL,
    status BOOLEAN,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_toggles_id ON toggles (id);

-- access_logs definition
CREATE TABLE access_logs (
    id SERIAL PRIMARY KEY,
    toggle_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_access_logs_toggle_id ON access_logs (toggle_id);