package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/outboxmsg"
)

type OutboxMessage struct {
	ID             int           `db:"id"`
	ChatID         int64         `db:"chat_id"`
	ReplyMessageID sql.NullInt64 `db:"reply_msg_id"`
	Text           string        `db:"msg_text"`
	Type           int           `db:"msg_type"`
	Payload        []byte        `db:"payload"`
	Button         jsonb.JSONB   `db:"buttons"`
	Status         string        `db:"status"`
	RetryCount     int           `db:"retry_count"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      sql.NullTime  `db:"updated_at"`
	VisibilityAt   time.Time     `db:"visibility_at"`
}

func intToNullInt64(value int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(value),
		Valid: value != 0,
	}
}

func ToDBOutboxMessage(msg msginfo.Message) (OutboxMessage, error) {
	jbButtons, err := jsonbFromButtonRows(msg.Buttons)
	if err != nil {
		return OutboxMessage{}, fmt.Errorf("jsonb from buttons: %w", err)
	}

	return OutboxMessage{
		ChatID:         msg.ChatID.Int64(),
		ReplyMessageID: intToNullInt64(msg.ReplyMsgID.Int()),
		Text:           msg.Text,
		Type:           msg.Type.Int(),
		Payload:        msg.Payload,
		Button:         jbButtons,
		VisibilityAt:   msg.VisibilityAt,
	}, nil
}

func jsonbFromButtonRows(buttons []button.ButtonRow) (jsonb.JSONB, error) {
	if buttons == nil {
		return jsonb.NewString("[]"), nil
	}

	jbButtons, err := jsonb.NewFromMarshaler(buttons)
	if err != nil {
		return jsonb.NewNull(), fmt.Errorf("jsonb marshaler: %w", err)
	}

	return jbButtons, nil
}

func ToOutboxMessage(msg OutboxMessage) (outboxmsg.Message, error) {
	var buttons []button.ButtonRow
	if err := jsonb.ConvertTo(msg.Button, &buttons); err != nil {
		return outboxmsg.Message{}, fmt.Errorf("convert jsonb to button rows: %w", err)
	}

	return outboxmsg.Message{
		ID:         msg.ID,
		RetryCount: msg.RetryCount,
		Message: msginfo.Message{
			ChatID:     msginfo.ChatIDFromInt64(msg.ChatID),
			ReplyMsgID: msginfo.MessageIDFromInt(int(msg.ReplyMessageID.Int64)),
			Text:       msg.Text,
			Type:       msginfo.MessageTypeFromInt(msg.Type),
			Payload:    msg.Payload,
			Buttons:    buttons,
		},
	}, nil
}

func ToOutboxMessages(dbMsgs []OutboxMessage) ([]outboxmsg.Message, error) {
	if len(dbMsgs) == 0 {
		return nil, nil
	}

	outboxMsgs := make([]outboxmsg.Message, 0, len(dbMsgs))

	for _, m := range dbMsgs {
		outboxMsg, err := ToOutboxMessage(m)
		if err != nil {
			return nil, fmt.Errorf("make outbox message: %w", err)
		}

		outboxMsgs = append(outboxMsgs, outboxMsg)
	}

	return outboxMsgs, nil
}
