package leavegame

import (
	"context"
	"fmt"
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

type GameLeaver interface {
	LeaveGame(
		ctx context.Context,
		playerID player.ID,
		gameID game.ID,
		completedAt time.Time,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type LeaveGame struct {
	playerProvider PlayerProvider
	gameLeaver     GameLeaver
	timeProvider   TimeProvider
}

func New(
	playerProvider PlayerProvider,
	gameLeaver GameLeaver,
	timeProvider TimeProvider,
) *LeaveGame {
	return &LeaveGame{
		playerProvider: playerProvider,
		gameLeaver:     gameLeaver,
		timeProvider:   timeProvider,
	}
}

func (s *LeaveGame) LeaveGame(ctx context.Context, chatID msginfo.ChatID) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("gat player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusIdle:
		return nil

	case player.GameStatusReadyForGame:
		return nil

	case player.GameStatusInGame:
		if err := s.gameLeaver.LeaveGame(
			ctx,
			currentPlayer.ID,
			game.IDFromString(currentPlayer.CurrentGameID),
			s.timeProvider.Now(),
		); err != nil {
			return fmt.Errorf("leave game: %w", err)
		}
	}

	return nil
}
