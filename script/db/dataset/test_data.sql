BEGIN;

-- INSERT INTO player(chat_id, display_name, created_at) VALUES(123, 'test1', NOW());
-- INSERT INTO player(chat_id, display_name, created_at, game_status) VALUES(124, 'test2', NOW(), 'ready_for_game');

INSERT INTO player(chat_id, display_name, created_at, game_status_changed_at, score)
SELECT 
    i,
    'player_' || i,
    NOW(),
    NOW(),
    500 + i
FROM generate_series(1, 10000) AS i;

COMMIT;