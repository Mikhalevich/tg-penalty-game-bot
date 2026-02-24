package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Game struct {
	ID              string      `db:"id"`
	CreatedAt       time.Time   `db:"created_at"`
	Status          string      `db:"game_status"`
	StatusChagnedAt time.Time   `db:"game_status_changed_at"`
	Payload         jsonb.JSONB `db:"payload"`
	PayloadVersion  int         `db:"payload_version"`
}

func (g Game) ToGame() (game.Game, error) {
	var gameState game.State

	if err := jsonb.ConvertTo(g.Payload, &gameState); err != nil {
		return game.Game{}, fmt.Errorf("convert payload to game state: %w", err)
	}

	return game.Game{
		ID:              game.IDFromString(g.ID),
		CreatedAt:       g.CreatedAt,
		Status:          game.GameStatus(g.Status),
		StatusChagnedAt: g.StatusChagnedAt,
		State:           gameState,
		StateVersion:    g.PayloadVersion,
	}, nil
}

func ToDBGame(domGame game.Game) (Game, error) {
	payload, err := jsonb.NewFromMarshaler(domGame.State)
	if err != nil {
		return Game{}, fmt.Errorf("marshal payload: %w", err)
	}

	return Game{
		ID:              domGame.ID.String(),
		CreatedAt:       domGame.CreatedAt,
		Status:          domGame.Status.String(),
		StatusChagnedAt: domGame.StatusChagnedAt,
		Payload:         payload,
		PayloadVersion:  domGame.StateVersion,
	}, nil
}

func ToDBGames(games []game.Game) ([]Game, error) {
	dbGames := make([]Game, 0, len(games))

	for _, g := range games {
		dbGame, err := ToDBGame(g)
		if err != nil {
			return nil, fmt.Errorf("convert to db game: %w", err)
		}

		dbGames = append(dbGames, dbGame)
	}

	return dbGames, nil
}
