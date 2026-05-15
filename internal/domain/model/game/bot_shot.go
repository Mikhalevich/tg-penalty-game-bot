package game

import (
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

var (
	//nolint:gochecknoglobals
	possibleExpectedShotSides = []ShotSide{
		ShotSideLeft,
		ShotSideMiddle,
		ShotSideRight,
	}
)

func (g *Game) BotShot(
	round int,
	completedAt time.Time,
) error {
	if !g.IsGameWithBot() {
		return perror.InvalidParam("game is not with bot")
	}

	botShot, playerShot := g.currentRoundShots(0)

	if playerShot.ActualSide == ShotSideNoShot {
		return perror.InvalidParam("player is not shot yet")
	}

	if botShot.ExpectedSide != ShotSideNoShot {
		return perror.AlreadyExists("shot already exist")
	}

	if botShot.Round != round {
		return perror.InvalidRound()
	}

	botShot.ExpectedSide = generateExpectedShotSide()
	botShot.ActualSide = calculateActualBotShotSide(
		botShot.Type,
		g.State.PlayerByID(0).Info,
		botShot.ExpectedSide,
		g.State.PlayerByID(playerShot.PlayerID).Info,
		playerShot.ActualSide,
	)

	botShot.CompletedAt = completedAt

	g.StateUpdatedAt = completedAt

	return nil
}

func calculateActualBotShotSide(
	shotType ShotType,
	botInfo PlayerInfo,
	botShotSide ShotSide,
	playerInfo PlayerInfo,
	playerShotSide ShotSide,
) ShotSide {
	if isMissPercentMatch(botInfo, shotType) {
		return ShotSideMiss
	}

	if playerShotSide == ShotSideMiss {
		return botShotSide
	}

	switch shotType {
	case ShotTypeAttack:
		return calculateActualAttackBotShotSide(botInfo, botShotSide, playerInfo, playerShotSide)

	case ShotTypeDefend:
		return calculateActualDefendBotShotSide(botInfo, botShotSide, playerInfo, playerShotSide)
	}

	return botShotSide
}

func calculateActualAttackBotShotSide(
	botInfo PlayerInfo, botShotSide ShotSide,
	playerInfo PlayerInfo, playerShotSide ShotSide,
) ShotSide {
	if isPercentMatch(botInfo.ForceGoalPercent) {
		return firstNotEqualSide(playerShotSide)
	}

	if isPercentMatch(playerInfo.ForceSavePercent) {
		return playerShotSide
	}

	return botShotSide
}

func calculateActualDefendBotShotSide(
	botInfo PlayerInfo, botShotSide ShotSide,
	playerInfo PlayerInfo, playerShotSide ShotSide,
) ShotSide {
	if isPercentMatch(botInfo.ForceSavePercent) {
		return playerShotSide
	}

	if isPercentMatch(playerInfo.ForceGoalPercent) {
		return firstNotEqualSide(playerShotSide)
	}

	return botShotSide
}

func firstNotEqualSide(side ShotSide) ShotSide {
	for _, s := range possibleExpectedShotSides {
		if s != side {
			return s
		}
	}

	return side
}

func generateExpectedShotSide() ShotSide {
	//nolint:gosec
	return possibleExpectedShotSides[rand.Int()%len(possibleExpectedShotSides)]
}

// currentRoundShots returns shot by player id for modification and shot another player for ready only access.
func (g *Game) currentRoundShots(id player.ID) (*Shot, Shot) {
	cRound := g.currentRoundPtr()

	if cRound.Attack.PlayerID == id {
		return &cRound.Attack, cRound.Defend
	}

	return &cRound.Defend, cRound.Attack
}
