package startgame

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	playersPerGame = 2
)

type Repository interface {
	InsertGames(ctx context.Context, games []game.Game) error
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
	GameNewRound(ctx context.Context, gameID game.ID, state game.State) error
	GameStart(
		ctx context.Context,
		gameType game.GameType,
		player1, player2 game.Player,
	) error
	GameLink(ctx context.Context, chatID msginfo.ChatID, gameID game.ID) error
}

type StartGame struct {
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
) *StartGame {
	return &StartGame{
		repo:                repo,
		transactor:          transactor,
		playerStatusChanger: playerStatusChanger,
		notifier:            notifier,
	}
}

func (s *StartGame) StartGames(ctx context.Context, games []game.Game) error {
	if len(games) == 0 {
		return errors.New("invalid games count")
	}

	if err := s.processStartGame(ctx, games); err != nil {
		var (
			playerIDs = make([]player.ID, 0, len(games)*playersPerGame)
			changedAt = games[0].CreatedAt
		)

		for _, currentGame := range games {
			playerIDs = append(playerIDs, currentGame.State.Player1.ID, currentGame.State.Player2.ID)
		}

		if err := s.playerStatusChanger.ChangePlayersGameStatus(
			ctx,
			playerIDs,
			player.GameStatusIdle,
			changedAt,
		); err != nil {
			return fmt.Errorf("change player to idle status: %w", err)
		}

		return fmt.Errorf("process start game: %w", err)
	}

	return nil
}

func (s *StartGame) processStartGame(ctx context.Context, games []game.Game) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.InsertGames(ctx, games); err != nil {
			return fmt.Errorf("insert games: %w", err)
		}

		for _, currentGame := range games {
			if err := s.sendNotificationForCreatedGame(ctx, currentGame); err != nil {
				return fmt.Errorf("send notification: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (s *StartGame) sendNotificationForCreatedGame(ctx context.Context, currentGame game.Game) error {
	switch currentGame.Status {
	case game.GameStatusPending:
		if err := s.notifier.GameLink(ctx, currentGame.State.Player1.ChatID, currentGame.ID); err != nil {
			return fmt.Errorf("game link: %w", err)
		}

	case game.GameStatusInProgress:
		if err := s.notifier.GameStart(ctx, currentGame.Type,
			currentGame.State.Player1, currentGame.State.Player2,
		); err != nil {
			return fmt.Errorf("start game: %w", err)
		}

		if err := s.notifier.GameNewRound(ctx, currentGame.ID, currentGame.State); err != nil {
			return fmt.Errorf("new round: %w", err)
		}

	case game.GameStatusCanceled, game.GameStatusCompleted:
		fallthrough

	default:
		return fmt.Errorf("invalid status: %s", currentGame.Status.String())
	}

	return nil
}
