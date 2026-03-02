package matchmaking

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (m *MatchMaking) ProcessReadyToGamePlayers(
	ctx context.Context,
	playersLimit int,
) error {
	var (
		games []game.Game
		err   error
	)

	if err := m.transactor.Transaction(ctx, func(ctx context.Context) error {
		games, err = m.transactionReadyToGamePlayers(ctx, playersLimit)
		if err != nil {
			return fmt.Errorf("process ready to game players: %w", err)
		}

		for _, gm := range games {
			if err := m.notifier.GameStage(ctx, gm); err != nil {
				return fmt.Errorf("send game stage: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (m *MatchMaking) transactionReadyToGamePlayers(
	ctx context.Context,
	playersLimit int,
) ([]game.Game, error) {
	players, err := m.repo.SelectReadyToGamePlayers(ctx, playersLimit)
	if err != nil {
		return nil, fmt.Errorf("select players: %w", err)
	}

	if len(players) <= 1 {
		return nil, nil
	}

	var (
		inGamePlayers = make([]player.Player, 0, len(players))
		//nolint:mnd
		games = make([]game.Game, 0, len(players)/2)
	)

	for i := 1; i < len(players); i += 2 {
		var (
			player1 = players[i-1]
			player2 = players[i]
		)

		newGame, err := m.gameCreator.CreateGame(ctx, player1, player2)
		if err != nil {
			return nil, fmt.Errorf("create game: %w", err)
		}

		games = append(games, newGame)

		inGamePlayers = appendInGamePlayers(inGamePlayers, player1, newGame.ID, newGame.CreatedAt)
		inGamePlayers = appendInGamePlayers(inGamePlayers, player2, newGame.ID, newGame.CreatedAt)
	}

	if err := m.insertGames(ctx, games); err != nil {
		return nil, fmt.Errorf("insert games: %w", err)
	}

	if err := m.setPlayersInGameStatus(ctx, inGamePlayers); err != nil {
		return nil, fmt.Errorf("change players game status to in game: %w", err)
	}

	return games, nil
}

func appendInGamePlayers(
	players []player.Player,
	plr player.Player,
	gameID game.ID,
	changedTime time.Time,
) []player.Player {
	plr.GameStatus = player.GameStatusInGame
	plr.GameStatusChangedAt = changedTime
	plr.CurrentGameID = gameID.String()

	return append(players, plr)
}

func (m *MatchMaking) setPlayersInGameStatus(ctx context.Context, players []player.Player) error {
	if len(players) == 0 {
		return nil
	}

	if err := m.repo.ChangeOrInsertPlayersGameStatus(
		ctx,
		players,
	); err != nil {
		return fmt.Errorf("repo change players game status to in game: %w", err)
	}

	return nil
}

func (m *MatchMaking) insertGames(ctx context.Context, games []game.Game) error {
	if len(games) == 0 {
		return nil
	}

	if err := m.repo.InsertGames(ctx, games); err != nil {
		return fmt.Errorf("repo insert games: %w", err)
	}

	return nil
}
