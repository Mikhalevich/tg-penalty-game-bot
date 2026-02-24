package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Shot struct {
	GameID      string       `db:"game_id"`
	PlayerID    int          `db:"player_id"`
	Round       int          `db:"round"`
	Type        string       `db:"shot_type"`
	CreatedAt   time.Time    `db:"created_at"`
	Side        string       `db:"side"`
	CompletedAt sql.NullTime `db:"completed_at"`
}

func (s Shot) ToShot() game.Shot {
	return game.Shot{
		GameID:      game.IDFromString(s.GameID),
		PlayerID:    player.IDFromInt(s.PlayerID),
		Round:       s.Round,
		Type:        game.ShotType(s.Type),
		CreatedAt:   s.CreatedAt,
		Side:        game.ShotSide(s.Side),
		CompletedAt: s.CompletedAt.Time,
	}
}

func ToDBShot(domShot game.Shot) Shot {
	return Shot{
		GameID:      domShot.GameID.String(),
		PlayerID:    domShot.PlayerID.Int(),
		Round:       domShot.Round,
		Type:        domShot.Type.String(),
		CreatedAt:   domShot.CreatedAt,
		Side:        domShot.Side.String(),
		CompletedAt: toNullTime(domShot.CompletedAt),
	}
}

func ToDBShots(domShots []game.Shot) []Shot {
	dbShots := make([]Shot, 0, len(domShots))

	for _, s := range domShots {
		dbShots = append(dbShots, ToDBShot(s))
	}

	return dbShots
}
