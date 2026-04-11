package shotstats

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	PlayerShotStatsByCurrentMonth(
		ctx context.Context,
		playerID player.ID,
		shotType game.ShotType,
	) ([]game.ShotSidePercent, error)
}

type ShotStats struct {
	repo Repository
}

func New(repo Repository) *ShotStats {
	return &ShotStats{
		repo: repo,
	}
}

func (ss *ShotStats) PlayerStats(
	ctx context.Context,
	playerID player.ID,
) (game.ShotStats, error) {
	attack, err := ss.repo.PlayerShotStatsByCurrentMonth(ctx, playerID, game.ShotTypeAttack)
	if err != nil {
		return game.ShotStats{}, fmt.Errorf("attack player stats: %w", err)
	}

	defend, err := ss.repo.PlayerShotStatsByCurrentMonth(ctx, playerID, game.ShotTypeDefend)
	if err != nil {
		return game.ShotStats{}, fmt.Errorf("defend player stats: %w", err)
	}

	return game.ShotStats{
		Attack: attack,
		Defend: defend,
	}, nil
}
