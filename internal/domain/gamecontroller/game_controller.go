package gamecontroller

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
	UpdateShot(ctx context.Context, shot game.Shot) error
	GetRoundShotsByGame(ctx context.Context, gameID game.ID, round int) ([]game.Shot, error)
	IsNoRowsUpdated(err error) bool
}

type PlayerController interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
}

type TimeProvider interface {
	Now() time.Time
}

type GameController struct {
	repo             Repository
	playerController PlayerController
	timeProvider     TimeProvider
}

func New(
	repo Repository,
	playerController PlayerController,
	timeProvier TimeProvider,
) *GameController {
	return &GameController{
		repo:             repo,
		playerController: playerController,
		timeProvider:     timeProvier,
	}
}
