package lbrefresher

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type HistoryPositions struct {
	Year      int
	Month     int
	Positions []player.Position
}

type Repository interface {
	RefreshScoreLeaderboard(ctx context.Context) error
	PlayerScorePositionsFrom(
		ctx context.Context,
		positionFrom int,
		limit int,
	) ([]player.Position, error)
	InsertHistoryScoreLeaderboard(ctx context.Context, positions HistoryPositions) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type SettingsProvider interface {
	GetLeaderboardMonth(ctx context.Context) (int, error)
	SetLeaderboardMonth(ctx context.Context, year, month int) error
}

type TimeProvider interface {
	Now() time.Time
}

type LeaderboardRefresher struct {
	repo             Repository
	transactor       Transactor
	settingsProvider SettingsProvider
	timeProvider     TimeProvider
}

func New(
	repo Repository,
	transactor Transactor,
	settingsProvider SettingsProvider,
	timeProvider TimeProvider,
) *LeaderboardRefresher {
	return &LeaderboardRefresher{
		repo:             repo,
		transactor:       transactor,
		settingsProvider: settingsProvider,
		timeProvider:     timeProvider,
	}
}
