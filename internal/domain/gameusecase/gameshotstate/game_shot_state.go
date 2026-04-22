package gameshotstate

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
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
	repo     Repository
	notifier Notifier
}

func New(
	repo Repository,
	notifier Notifier,
) *GameShotState {
	return &GameShotState{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *GameShotState) ShotState(ctx context.Context, currentPlayer player.Player) error {
	if currentPlayer.GameStatus != player.GameStatusInGame {
		return errors.New("player not in game")
	}

	currentGame, err := s.repo.GetGame(ctx, game.IDFromString(currentPlayer.CurrentGameID))
	if err != nil {
		return fmt.Errorf("get game by id: %w", err)
	}

	if currentGame.IsPending() {
		return perror.AlreadyInGame()
	}

	if !currentGame.IsInProgress() {
		return fmt.Errorf("game is not in progress, game_id: %q state: %q",
			currentGame.ID.String(), currentGame.Status.String())
	}

	if err := s.sendShotStateNotification(
		ctx,
		currentPlayer,
		currentGame.ID,
		currentGame.State.CurrentRound(),
		currentGame.State.CurrentRoundNumber(),
	); err != nil {
		return fmt.Errorf("send shot state notification: %w", err)
	}

	return nil
}

func (s *GameShotState) sendShotStateNotification(
	ctx context.Context,
	currentPlayer player.Player,
	gameID game.ID,
	round game.Round,
	roundNumber int,
) error {
	if round.IsCompleted() {
		return fmt.Errorf("last raund is completed, game_id: %q round: %d",
			gameID.String(), roundNumber)
	}

	var shotType game.ShotType
	switch currentPlayer.ID {
	case round.Attack.PlayerID:
		shotType = game.ShotTypeAttack

	case round.Defend.PlayerID:
		shotType = game.ShotTypeDefend

	default:
		return fmt.Errorf("invalid player for game, player_id: %d game_id: %q",
			currentPlayer.ID.Int(), gameID.String())
	}

	if err := s.notifier.GameShotState(
		ctx,
		currentPlayer.ChatID,
		gameID,
		shotType,
		roundNumber,
	); err != nil {
		return fmt.Errorf("game shot state notification: %w", err)
	}

	return nil
}
