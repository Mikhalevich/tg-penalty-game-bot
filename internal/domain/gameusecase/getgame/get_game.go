package getgame

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
}

type GetGame struct {
	repo Repository
}

func New(
	repo Repository,
) *GetGame {
	return &GetGame{
		repo: repo,
	}
}

func (s *GetGame) GetGameByID(ctx context.Context, gameID game.ID) (game.Game, error) {
	currentGame, err := s.repo.GetGame(ctx, gameID)
	if err != nil {
		return game.Game{}, fmt.Errorf("repo get game: %w", err)
	}

	return currentGame, nil
}
