package game

import (
	"math/rand"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (g *Game) PlayerShot(shot Shot) error {
	inGameShot := g.shotByPlayerID(shot.PlayerID)

	if inGameShot.ExpectedSide != ShotSideNoShot {
		return perror.AlreadyExists("shot already exist")
	}

	if inGameShot.PlayerID != shot.PlayerID {
		return perror.InvalidPlayer()
	}

	if inGameShot.Round != shot.Round {
		return perror.InvalidRound()
	}

	inGameShot.ExpectedSide = shot.ExpectedSide
	inGameShot.ActualSide = calculateActualShotSide(
		g.State.PlayerByID(shot.PlayerID).Info,
		shot.Type,
		shot.ExpectedSide,
	)
	inGameShot.CompletedAt = shot.CompletedAt

	g.StateUpdatedAt = shot.CompletedAt

	return nil
}

func isPercentMatch(percent int) bool {
	if percent == 0 {
		return false
	}

	//nolint:gosec
	return rand.Int()%percent100 < percent
}

func calculateActualShotSide(info PlayerInfo, shotType ShotType, side ShotSide) ShotSide {
	if isMissPercentMatch(info, shotType) {
		return ShotSideMiss
	}

	return side
}

func isMissPercentMatch(info PlayerInfo, shotType ShotType) bool {
	missPercent := 0
	switch shotType {
	case ShotTypeAttack:
		missPercent = info.AttackMissPercent

	case ShotTypeDefend:
		missPercent = info.DefendMissPercent
	}

	return isPercentMatch(missPercent)
}
