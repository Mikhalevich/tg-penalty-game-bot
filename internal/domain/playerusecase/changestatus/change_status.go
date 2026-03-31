package changestatus

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	ChangePlayersGameStatus(
		ctx context.Context,
		playerIDs []player.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error
	SetPlayerIdleStatusWithScore(
		ctx context.Context,
		playerID player.ID,
		scoreDelta int,
		changedAt time.Time,
	) error
}

type ChangeStatus struct {
	repo Repository
}

func New(
	repo Repository,
) *ChangeStatus {
	return &ChangeStatus{
		repo: repo,
	}
}
