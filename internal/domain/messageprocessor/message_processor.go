package messageprocessor

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/shotimage"
)

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg msginfo.SenderMessage,
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

type ShotImageProvider interface {
	Image(ctx context.Context, shot shotimage.ShotImage) ([]byte, error)
}

type MessageProcessor struct {
	sender            Sender
	escaper           MarkdownEscaper
	buttonRepository  ButtonRepository
	shotImageProvider ShotImageProvider
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
	shotImageProvider ShotImageProvider,
) *MessageProcessor {
	return &MessageProcessor{
		sender:            sender,
		escaper:           escaper,
		buttonRepository:  buttonRepository,
		shotImageProvider: shotImageProvider,
	}
}
