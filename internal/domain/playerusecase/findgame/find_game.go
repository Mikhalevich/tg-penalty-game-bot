package findgame

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	ChangePlayerGameStatus(
		ctx context.Context,
		playerID player.ID,
		gameID game.ID,
		status player.GameStatus,
		changedTime time.Time,
		prevoiusStatuses ...player.GameStatus,
	) error
}

type PlayerProvider interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
}

type TimeProvider interface {
	Now() time.Time
}

type ShotStater interface {
	ShotState(ctx context.Context, currentPlayer player.Player) error
}

type FindGame struct {
	repo           Repository
	playerProvider PlayerProvider
	timeProvider   TimeProvider
	shotStater     ShotStater
}

func New(
	repo Repository,
	playerProvider PlayerProvider,
	timeProvider TimeProvider,
	shotStater ShotStater,
) *FindGame {
	return &FindGame{
		repo:           repo,
		playerProvider: playerProvider,
		timeProvider:   timeProvider,
		shotStater:     shotStater,
	}
}

func (fg *FindGame) FindGame(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := fg.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusInGame:
		if err := fg.shotStater.ShotState(ctx, currentPlayer); err != nil {
			return fmt.Errorf("shot state: %w", err)
		}

		return nil

	case player.GameStatusReadyForGame:
		return nil

	case player.GameStatusIdle:
	}

	if err := fg.repo.ChangePlayerGameStatus(
		ctx,
		currentPlayer.ID,
		"",
		player.GameStatusReadyForGame,
		fg.timeProvider.Now(),
		player.GameStatusIdle,
	); err != nil {
		return fmt.Errorf("change player game status: %w", err)
	}

	return nil
}
