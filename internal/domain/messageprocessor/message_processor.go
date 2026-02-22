package messageprocessor

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

type MessageType int

const (
	MessageTypePlain MessageType = iota + 1
	MessageTypeMarkdown
	MessageTypePNG
)

func (mt MessageType) Int() int {
	return int(mt)
}

func MessageTypeFromInt(t int) MessageType {
	return MessageType(t)
}

type Message struct {
	ChatID     msginfo.ChatID
	ReplyMsgID msginfo.MessageID
	Text       string
	Type       MessageType
	Payload    []byte
	Buttons    []button.ButtonRow
}

type SenderMessage struct {
	ChatID     msginfo.ChatID
	ReplyMsgID msginfo.MessageID
	Text       string
	Type       MessageType
	Payload    []byte
	Buttons    []button.InlineKeyboardButtonRow
}

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg SenderMessage,
	) error
	EditText(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		text string,
		rows ...button.InlineKeyboardButtonRow,
	) error
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type ButtonRepository interface {
	SetButton(ctx context.Context, btn button.Button) error
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) error

	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	IsNotFoundError(err error) bool
}

type MessageProcessor struct {
	sender           Sender
	escaper          MarkdownEscaper
	buttonRepository ButtonRepository
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
) *MessageProcessor {
	return &MessageProcessor{
		sender:           sender,
		escaper:          escaper,
		buttonRepository: buttonRepository,
	}
}
