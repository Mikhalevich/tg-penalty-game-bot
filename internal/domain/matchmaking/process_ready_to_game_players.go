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
	if err := m.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := m.transactionReadyToGamePlayers(ctx, playersLimit); err != nil {
			return fmt.Errorf("process ready to game players: %w", err)
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
) error {
	players, err := m.repo.SelectReadyToGamePlayers(ctx, playersLimit)
	if err != nil {
		return fmt.Errorf("select players: %w", err)
	}

	if len(players) <= 1 {
		return nil
	}

	var (
		inGamePlayers = make([]player.Player, 0, len(players))
		//nolint:mnd
		games = make([]game.Game, 0, len(players)/2)
		now   = m.timeProvider.Now()
	)

	for i := 1; i < len(players); i += 2 {
		var (
			player1 = players[i-1]
			player2 = players[i]
		)

		newGame := game.CreateGame(ctx, game.GameTypeRating, player1, player2, now)

		games = append(games, newGame)

		inGamePlayers = appendInGamePlayers(inGamePlayers, player1, newGame.ID, newGame.CreatedAt)
		inGamePlayers = appendInGamePlayers(inGamePlayers, player2, newGame.ID, newGame.CreatedAt)
	}

	if len(games) > 0 {
		if err := m.gameRunner.StartGames(ctx, games); err != nil {
			return fmt.Errorf("insert games: %w", err)
		}

		if err := m.repo.ChangeOrInsertPlayersGameStatus(ctx, inGamePlayers); err != nil {
			return fmt.Errorf("change players game status to in game: %w", err)
		}
	}

	return nil
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
