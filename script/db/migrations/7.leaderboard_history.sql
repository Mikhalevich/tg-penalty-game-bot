-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE leaderboard_score_history(
    player_id INTEGER NOT NULL,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    display_name TEXT NOT NULL,
    score INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    position INTEGER NOT NULL,

    CONSTRAINT pk_leaderboard_score_history PRIMARY KEY (player_id, year, month)
);

CREATE INDEX idx_leaderboard_score_history_year_month ON leaderboard_score_history(year, month);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE leaderboard_score_history;
