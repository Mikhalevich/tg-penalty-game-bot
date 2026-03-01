package matchmaking

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	SelectReadyToGamePlayers(ctx context.Context, limit int) ([]player.Player, error)
	InsertGames(ctx context.Context, games []game.Game) error
	ChangeOrInsertPlayersGameStatus(
		ctx context.Context,
		players []player.Player,
	) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type GameCreator interface {
	CreateGame(ctx context.Context, players []player.Player) (game.Game, error)
}

type Notifier interface {
	GameStage(ctx context.Context, gm game.Game) error
}

type MatchMaking struct {
	repo        Repository
	transactor  Transactor
	gameCreator GameCreator
	notifier    Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	gameCreator GameCreator,
	notifier Notifier,
) *MatchMaking {
	return &MatchMaking{
		repo:        repo,
		transactor:  transactor,
		gameCreator: gameCreator,
		notifier:    notifier,
	}
}
