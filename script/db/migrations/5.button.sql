-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE button(
    id TEXT PRIMARY KEY,
    caption TEXT NOT NULL,
    operation TEXT NOT NULL,
    is_delete_message BOOLEAN NOT NULL,
    payload BYTEA
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE button;
