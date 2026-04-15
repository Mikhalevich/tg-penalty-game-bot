-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE game_type AS ENUM (
    'rating',
    'friendly'
);

CREATE TYPE game_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'canceled'
);

CREATE TABLE game(
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    game_type game_type NOT NULL,
    game_status game_status NOT NULL,
    payload JSONB NOT NULL,
    payload_version INTEGER NOT NULL DEFAULT 0,
    payload_updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX game_raring_in_progress_idx ON game(game_type, game_status, payload_updated_at) WHERE game_type = 'rating' AND game_status = 'in_progress';

CREATE TYPE shot_type AS ENUM (
    'defend',
    'attack'
);

CREATE TYPE shot_side AS ENUM (
    'no_shot',
    'left',
    'right',
    'middle',
    'miss'
);

CREATE TABLE shot(
    game_id UUID NOT NULL,
    player_id INTEGER NOT NULL,
    round INTEGER NOT NULL,
    shot_type shot_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    expected_side shot_side NOT NULL DEFAULT 'no_shot',
    actual_side shot_side NOT NULL DEFAULT 'no_shot',
    completed_at TIMESTAMPTZ,

    CONSTRAINT pk_shots PRIMARY KEY (game_id, player_id, round),
    CONSTRAINT fk_shots_player FOREIGN KEY(player_id) REFERENCES player(id),
    CONSTRAINT fk_shots_game FOREIGN KEY(game_id) REFERENCES game(id)
);

CREATE INDEX shot_player_id_created_at_idx ON shot(player_id, shot_type, created_at);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE shot;
DROP TYPE shot_side;
DROP TYPE shot_type;
DROP TABLE game;
DROP TYPE game_status;
DROP TYPE game_type;
