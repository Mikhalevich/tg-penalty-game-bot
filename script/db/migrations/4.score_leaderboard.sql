-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE MATERIALIZED VIEW score_leaderboard AS
    SELECT
        id AS player_id,
        display_name,
        score,
        game_status_changed_at AS updated_at,
        ROW_NUMBER() OVER (ORDER BY score DESC, game_status_changed_at) AS position
    FROM
        player
    WHERE
        game_status_changed_at >= DATE_TRUNC('month', CURRENT_DATE)
    ORDER BY
        position;

CREATE INDEX score_leaderboard_user_id_idx ON score_leaderboard(player_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE score_leaderboard;
