-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE settings(
    id TEXT PRIMARY KEY,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL,
    payload JSONB NOT NULL,
    payload_updated_at TIMESTAMPTZ NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE settings;
