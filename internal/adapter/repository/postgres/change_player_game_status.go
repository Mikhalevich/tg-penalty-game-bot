package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (p *Postgres) ChangePlayerGameStatus(
	ctx context.Context,
	playerID player.ID,
	gameID game.GameID,
	status player.GameStatus,
	changedTime time.Time,
	previousStatuses ...player.GameStatus,
) error {
	query, args, err := makeChangeGameStatusQuery(playerID, gameID, status, changedTime, previousStatuses)
	if err != nil {
		return fmt.Errorf("make change game status query: %w", err)
	}

	trx := p.transactor.ExtContext(ctx)

	res, err := trx.ExecContext(ctx, trx.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func makeChangeGameStatusQuery(
	playerID player.ID,
	gameID game.GameID,
	status player.GameStatus,
	changedTime time.Time,
	previousStatuses []player.GameStatus,
) (string, []any, error) {
	queryTemplate := `
		UPDATE player SET
			game_status = :game_status,
			game_status_changed_at = :game_status_changed_at,
			current_game_id = :current_game_id
		WHERE
			id = :id
			%s
	`

	if len(previousStatuses) == 0 {
		return bindChangeGameStatusQuery(
			fmt.Sprintf(queryTemplate, ""),
			playerID,
			gameID,
			status,
			changedTime,
		)
	}

	query, args, err := bindChangeGameStatusQuery(
		fmt.Sprintf(queryTemplate, "AND game_status IN(?)"),
		playerID,
		gameID,
		status,
		changedTime,
	)

	if err != nil {
		return "", nil, fmt.Errorf("bind query with in statement: %w", err)
	}

	args = append(args, previousStatuses)

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, fmt.Errorf("make in statement: %w", err)
	}

	return query, args, nil
}

func bindChangeGameStatusQuery(
	query string,
	playerID player.ID,
	gameID game.GameID,
	status player.GameStatus,
	changedTime time.Time,
) (string, []any, error) {
	//nolint:wrapcheck
	return sqlx.Named(
		query,
		map[string]any{
			"game_status":            status,
			"game_status_changed_at": changedTime,
			"id":                     playerID.Int(),
			"current_game_id":        gameID.String(),
		})
}
