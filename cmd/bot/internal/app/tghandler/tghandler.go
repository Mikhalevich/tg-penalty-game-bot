package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type PlayerController interface {
	Welcome(ctx context.Context, chatID msginfo.ChatID) error
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID, fullName, userName string) error
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, msgID msginfo.MessageID, displayName string) error
	SetPlayerReadyToGame(ctx context.Context, chatID msginfo.ChatID) error
}

type GameController interface {
	Shot(
		ctx context.Context,
		chatID msginfo.ChatID,
		msgID msginfo.MessageID,
		gameID game.ID,
		round int,
		side game.ShotSide,
	) error
	StartGameWithBot(ctx context.Context, chatID msginfo.ChatID) error
}

type TGHandler struct {
	buttonProvider   ButtonProvider
	playerController PlayerController
	gameController   GameController
}

func New(
	buttonProvider ButtonProvider,
	playerController PlayerController,
	gameController GameController,
) *TGHandler {
	return &TGHandler{
		buttonProvider:   buttonProvider,
		playerController: playerController,
		gameController:   gameController,
	}
}
