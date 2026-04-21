package leavegame

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

type Notifier interface {
	StopSearchGame(ctx context.Context, chatID msginfo.ChatID) error
}

type LeaveGame struct {
	repo           Repository
	playerProvider PlayerProvider
	gameLeaver     GameLeaver
	timeProvider   TimeProvider
	notifier       Notifier
}

func New(
	repo Repository,
	playerProvider PlayerProvider,
	gameLeaver GameLeaver,
	timeProvider TimeProvider,
	notifier Notifier,
) *LeaveGame {
	return &LeaveGame{
		repo:           repo,
		playerProvider: playerProvider,
		gameLeaver:     gameLeaver,
		timeProvider:   timeProvider,
		notifier:       notifier,
	}
}

func (s *LeaveGame) LeaveGame(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := s.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("gat player by chat_id: %w", err)
	}

	switch currentPlayer.GameStatus {
	case player.GameStatusIdle:

	case player.GameStatusReadyForGame:
		if err := s.repo.ChangePlayerGameStatus(
			ctx,
			currentPlayer.ID,
			"",
			player.GameStatusIdle,
			s.timeProvider.Now(),
			player.GameStatusReadyForGame,
		); err != nil {
			return fmt.Errorf("change player status: %w", err)
		}

		if err := s.notifier.StopSearchGame(ctx, currentPlayer.ChatID); err != nil {
			return fmt.Errorf("stop search game: %w", err)
		}

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
