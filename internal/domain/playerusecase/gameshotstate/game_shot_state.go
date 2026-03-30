package gameshotstate

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type GameProvider interface {
	GetGameByID(ctx context.Context, gameID game.ID) (game.Game, error)
}

type Notifier interface {
	GameShotState(
		ctx context.Context,
		chatID msginfo.ChatID,
		gameID game.ID,
		shotType game.ShotType,
		roundNumber int,
	) error
}

type GameShotState struct {
	gameProvider GameProvider
	notifier     Notifier
}

func New(
	gameProvider GameProvider,
	notifier Notifier,
) *GameShotState {
	return &GameShotState{
		gameProvider: gameProvider,
		notifier:     notifier,
	}
}

func (s *GameShotState) ShotState(ctx context.Context, currentPlayer player.Player) error {
	if currentPlayer.GameStatus != player.GameStatusInGame {
		return errors.New("player not in game")
	}

	currentGame, err := s.gameProvider.GetGameByID(ctx, game.IDFromString(currentPlayer.CurrentGameID))
	if err != nil {
		return fmt.Errorf("get game by id: %w", err)
	}

	if !currentGame.IsInProgress() {
		return fmt.Errorf("game is not in progress, game_id: %q state: %q",
			currentGame.ID.String(), currentGame.Status.String())
	}

	var (
		currentRound = currentGame.State.CurrentRound()
		roundNumber  = currentGame.State.CurrentRoundNumber()
	)

	if currentRound.IsCompleted() {
		return fmt.Errorf("last raund is completed, game_id: %q round: %d",
			currentGame.ID.String(), roundNumber)
	}

	var shotType game.ShotType
	switch currentPlayer.ID {
	case currentRound.Attack.PlayerID:
		shotType = game.ShotTypeAttack

	case currentRound.Defend.PlayerID:
		shotType = game.ShotTypeDefend

	default:
		return fmt.Errorf("invalid player for game, player_id: %d game_id: %q",
			currentPlayer.ID.Int(), currentGame.ID.String())
	}

	if err := s.notifier.GameShotState(
		ctx,
		currentPlayer.ChatID,
		currentGame.ID,
		shotType,
		roundNumber,
	); err != nil {
		return fmt.Errorf("game shot state notification: %w", err)
	}

	return nil
}
