-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE game_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'canceled'
);

CREATE TABLE game(
    id UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    game_status game_status NOT NULL,
    game_status_changed_at TIMESTAMPTZ NOT NULL,
    payload JSONB NOT NULL,
    payload_version INTEGER NOT NULL DEFAULT 0
);

CREATE TYPE shot_type AS ENUM (
    'defend',
    'attack'
);

CREATE TYPE shot_side AS ENUM (
    'no_shot',
    'left',
    'right',
    'middle'
);

CREATE TABLE shot(
    game_id UUID NOT NULL,
    player_id INTEGER NOT NULL,
    round INTEGER NOT NULL,
    shot_type shot_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    side shot_side NOT NULL DEFAULT 'no_shot',
    completed_at TIMESTAMPTZ,

    CONSTRAINT fk_shots_player FOREIGN KEY(player_id) REFERENCES player(id),
    CONSTRAINT fk_shots_game FOREIGN KEY(game_id) REFERENCES game(id),
    CONSTRAINT pk_shots PRIMARY KEY (game_id, player_id, round)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE shot;
DROP TYPE shot_side;
DROP TYPE shot_type;
DROP TABLE game;
DROP TYPE game_status;
