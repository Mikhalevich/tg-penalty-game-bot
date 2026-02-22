BEGIN;

INSERT INTO player(chat_id, display_name, created_at) VALUES(123, 'test1', NOW());
INSERT INTO player(chat_id, display_name, created_at, game_status) VALUES(124, 'test2', NOW(), 'ready_for_game');

COMMIT;