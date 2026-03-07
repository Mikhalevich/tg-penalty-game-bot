package msginfo

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

type MessageID int

func (m MessageID) Int() int {
	return int(m)
}

func MessageIDFromInt(id int) MessageID {
	return MessageID(id)
}

type ChatID int64

func (id ChatID) Int64() int64 {
	return int64(id)
}

func ChatIDFromInt64(id int64) ChatID {
	return ChatID(id)
}

type MessageType int

const (
	MessageTypePlain MessageType = iota + 1
	MessageTypeMarkdown
	MessageTypePNG
	MessageTypeShotImage
)

func (mt MessageType) Int() int {
	return int(mt)
}

func MessageTypeFromInt(t int) MessageType {
	return MessageType(t)
}

type Message struct {
	ChatID     ChatID
	ReplyMsgID MessageID
	Text       string
	Type       MessageType
	Payload    []byte
	Buttons    []button.ButtonRow
}

type SenderMessage struct {
	ChatID     ChatID
	ReplyMsgID MessageID
	Text       string
	Type       MessageType
	Payload    []byte
	Buttons    []button.InlineKeyboardButtonRow
}
