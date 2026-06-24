-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE tournament_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'canceled'
);

CREATE TABLE tournament(
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    tournament_status tournament_status NOT NULL,
    payload JSONB NOT NULL,
    payload_version INTEGER NOT NULL DEFAULT 0,
    payload_updated_at TIMESTAMPTZ NOT NULL
);

ALTER TYPE game_type ADD VALUE 'tournament';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE tournament;
DROP TYPE tournament_status;
