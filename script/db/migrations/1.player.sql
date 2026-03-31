-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE player_game_status AS ENUM (
    'idle',
    'ready_for_game',
    'in_game'
);

CREATE TABLE player(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    is_change_name_triggered BOOLEAN NOT NULL DEFAULT FALSE,
    name_changed_at TIMESTAMPTZ,
    game_status player_game_status NOT NULL DEFAULT 'idle',
    game_status_changed_at TIMESTAMPTZ,
    current_game_id TEXT,
    score INTEGER NOT NULL
);

CREATE UNIQUE INDEX player_chat_id_u_idx ON player(chat_id);
CREATE UNIQUE INDEX player_display_name_u_idx ON player(display_name);
CREATE INDEX player_ready_for_game_idx ON player(game_status, game_status_changed_at) WHERE game_status = 'ready_for_game';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE player;
DROP TYPE player_game_status;
