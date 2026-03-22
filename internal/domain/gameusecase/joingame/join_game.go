package joingame

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Repository interface {
	GetGame(ctx context.Context, gameID game.ID) (game.Game, error)
	UpdateGame(ctx context.Context, game game.Game) error
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerStatusChanger interface {
	ChangePlayersGameStatus(
		ctx context.Context,
		playerIDs []player.ID,
		status player.GameStatus,
		changedAt time.Time,
	) error
}

type Notifier interface {
	GameStart(ctx context.Context, player1, player2 game.Player) error
	GameNewRound(ctx context.Context, gameID game.ID, state game.State) error
}

type JoinGame struct {
	repo                Repository
	transactor          Transactor
	playerStatusChanger PlayerStatusChanger
	notifier            Notifier
}

func New(
	repo Repository,
	transactor Transactor,
	playerStatusChanger PlayerStatusChanger,
	notifier Notifier,
) *JoinGame {
	return &JoinGame{
		repo:                repo,
		transactor:          transactor,
		playerStatusChanger: playerStatusChanger,
		notifier:            notifier,
	}
}

func (j *JoinGame) JoinGame(
	ctx context.Context,
	gameID game.ID,
	plr game.Player,
	joinedAt time.Time,
) error {
	if err := j.processJoin(ctx, gameID, plr, joinedAt); err != nil {
		if err := j.playerStatusChanger.ChangePlayersGameStatus(
			ctx,
			[]player.ID{plr.ID},
			player.GameStatusIdle,
			joinedAt,
		); err != nil {
			return fmt.Errorf("change player status to idle: %w", err)
		}

		return fmt.Errorf("process join game: %w", err)
	}

	return nil
}

func (j *JoinGame) processJoin(
	ctx context.Context,
	gameID game.ID,
	plr game.Player,
	joinedAt time.Time,
) error {
	if err := j.transactor.Transaction(ctx, func(ctx context.Context) error {
		currentGame, err := j.repo.GetGame(ctx, gameID)
		if err != nil {
			return fmt.Errorf("get game: %w", err)
		}

		if err := currentGame.JoinPlayerAndStartGame(plr, joinedAt); err != nil {
			return fmt.Errorf("join player: %w", err)
		}

		if err := j.notifier.GameStart(ctx, currentGame.State.Player1, currentGame.State.Player2); err != nil {
			return fmt.Errorf("start game: %w", err)
		}

		if err := j.notifier.GameNewRound(ctx, currentGame.ID, currentGame.State); err != nil {
			return fmt.Errorf("new round: %w", err)
		}

		if err := j.repo.UpdateGame(ctx, currentGame); err != nil {
			return fmt.Errorf("update game: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
