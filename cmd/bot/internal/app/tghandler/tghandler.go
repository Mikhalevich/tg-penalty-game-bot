package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type PlayerController interface {
	Welcome(ctx context.Context, chatID msginfo.ChatID) error
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID) error
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, displayName string) error
}

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type TGHandler struct {
	playerController PlayerController
	messageProcessor ButtonProvider
}

func New(
	playerController PlayerController,
	messageProcessor ButtonProvider,
) *TGHandler {
	return &TGHandler{
		playerController: playerController,
		messageProcessor: messageProcessor,
	}
}
