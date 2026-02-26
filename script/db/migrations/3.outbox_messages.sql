-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE outbox_messages(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    reply_msg_id BIGINT,
    msg_text TEXT NOT NULL,
    msg_type INTEGER NOT NULL,
    payload BYTEA,
    buttons JSONB NOT NULL,
    is_dispatched BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    dispatched_at TIMESTAMPTZ
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE outbox_messages;
