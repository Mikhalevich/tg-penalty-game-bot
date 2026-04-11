package model

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type ShotSidePercent struct {
	Side    string  `db:"side"`
	Percent float32 `db:"percent"`
}

func (ssp ShotSidePercent) ToDomPercent() game.ShotSidePercent {
	return game.ShotSidePercent{
		Side:    game.ShotSide(ssp.Side),
		Percent: ssp.Percent,
	}
}

func ToDomShotSidePercents(dbPercents []ShotSidePercent) []game.ShotSidePercent {
	if len(dbPercents) == 0 {
		return nil
	}

	domPercents := make([]game.ShotSidePercent, 0, len(dbPercents))

	for _, percent := range dbPercents {
		domPercents = append(domPercents, percent.ToDomPercent())
	}

	return domPercents
}
