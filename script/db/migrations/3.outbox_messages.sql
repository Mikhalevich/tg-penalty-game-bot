-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE outbox_message_status AS ENUM (
    'pending',
    'dispatched',
    'canceled'
);

CREATE TABLE outbox_messages(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    reply_msg_id BIGINT,
    msg_text TEXT NOT NULL,
    msg_type INTEGER NOT NULL,
    payload BYTEA,
    buttons JSONB NOT NULL,
    status outbox_message_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    updated_at TIMESTAMPTZ,
    visibility_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX outbox_messages_pending_idx ON outbox_messages(status, visibility_at) WHERE status = 'pending';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE outbox_messages;
DROP TYPE outbox_message_status;
