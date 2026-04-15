package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Shot struct {
	GameID       string       `db:"game_id"`
	PlayerID     int          `db:"player_id"`
	Round        int          `db:"round"`
	Type         string       `db:"shot_type"`
	CreatedAt    time.Time    `db:"created_at"`
	ExpectedSide string       `db:"expected_side"`
	ActualSide   string       `db:"actual_side"`
	CompletedAt  sql.NullTime `db:"completed_at"`
}

func (s Shot) ToShot() game.Shot {
	return game.Shot{
		GameID:       game.IDFromString(s.GameID),
		PlayerID:     player.IDFromInt(s.PlayerID),
		Round:        s.Round,
		Type:         game.ShotType(s.Type),
		CreatedAt:    s.CreatedAt,
		ExpectedSide: game.ShotSide(s.ExpectedSide),
		ActualSide:   game.ShotSide(s.ActualSide),
		CompletedAt:  s.CompletedAt.Time,
	}
}

func ToShots(dbShots []Shot) []game.Shot {
	domShots := make([]game.Shot, 0, len(dbShots))

	for _, s := range dbShots {
		domShots = append(domShots, s.ToShot())
	}

	return domShots
}

func ToDBShot(domShot game.Shot) Shot {
	return Shot{
		GameID:       domShot.GameID.String(),
		PlayerID:     domShot.PlayerID.Int(),
		Round:        domShot.Round,
		Type:         domShot.Type.String(),
		CreatedAt:    domShot.CreatedAt,
		ExpectedSide: domShot.ExpectedSide.String(),
		ActualSide:   domShot.ActualSide.String(),
		CompletedAt:  toNullTime(domShot.CompletedAt),
	}
}

func ToDBShots(domShots []game.Shot) []Shot {
	dbShots := make([]Shot, 0, len(domShots))

	for _, s := range domShots {
		dbShots = append(dbShots, ToDBShot(s))
	}

	return dbShots
}
