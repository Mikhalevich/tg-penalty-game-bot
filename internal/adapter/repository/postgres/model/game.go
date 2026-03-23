package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Game struct {
	ID               string      `db:"id"`
	CreatedAt        time.Time   `db:"created_at"`
	Type             string      `db:"game_type"`
	Status           string      `db:"game_status"`
	Payload          jsonb.JSONB `db:"payload"`
	PayloadVersion   int         `db:"payload_version"`
	PayloadUpdatedAt time.Time   `db:"payload_updated_at"`
}

func (g Game) ToGame() (game.Game, error) {
	var gameState game.State

	if err := jsonb.ConvertTo(g.Payload, &gameState); err != nil {
		return game.Game{}, fmt.Errorf("convert payload to game state: %w", err)
	}

	return game.Game{
		ID:             game.IDFromString(g.ID),
		CreatedAt:      g.CreatedAt,
		Type:           game.GameType(g.Type),
		Status:         game.GameStatus(g.Status),
		State:          gameState,
		StateVersion:   g.PayloadVersion,
		StateUpdatedAt: g.PayloadUpdatedAt,
	}, nil
}

func ToDBGame(domGame game.Game) (Game, error) {
	payload, err := jsonb.NewFromMarshaler(domGame.State)
	if err != nil {
		return Game{}, fmt.Errorf("marshal payload: %w", err)
	}

	return Game{
		ID:               domGame.ID.String(),
		CreatedAt:        domGame.CreatedAt,
		Type:             domGame.Type.String(),
		Status:           domGame.Status.String(),
		Payload:          payload,
		PayloadVersion:   domGame.StateVersion,
		PayloadUpdatedAt: domGame.StateUpdatedAt,
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
