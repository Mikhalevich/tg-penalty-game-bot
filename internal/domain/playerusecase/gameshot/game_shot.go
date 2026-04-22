package gameshot

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type PlayerProvider interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, error)
}

type GameController interface {
	Shot(ctx context.Context, shot game.Shot) error
}

type TimeProvider interface {
	Now() time.Time
}

type GameShot struct {
	playerProvider PlayerProvider
	gameController GameController
	timeProvider   TimeProvider
}

func New(
	playerProvider PlayerProvider,
	gameController GameController,
	timeProvider TimeProvider,
) *GameShot {
	return &GameShot{
		playerProvider: playerProvider,
		gameController: gameController,
		timeProvider:   timeProvider,
	}
}
