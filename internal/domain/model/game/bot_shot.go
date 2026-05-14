package game

import (
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

var (
	//nolint:gochecknoglobals
	possibleBotShotSides = []ShotSide{
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

	botShot, playerShot := g.playersShots(0)

	if playerShot.ActualSide == ShotSideNoShot {
		return perror.InvalidParam("player is not shot yet")
	}

	if botShot.ExpectedSide != ShotSideNoShot {
		return perror.AlreadyExists("shot already exist")
	}

	if botShot.Round != round {
		return perror.InvalidRound()
	}

	botShot.ExpectedSide = generateBotSide()
	botShot.ActualSide = calculateActualBotShotSide(
		g.State.PlayerByID(0).Info,
		botShot.Type,
		botShot.ExpectedSide,
		g.State.PlayerByID(playerShot.PlayerID).Info,
		playerShot.ActualSide,
	)

	botShot.CompletedAt = completedAt

	g.StateUpdatedAt = completedAt

	return nil
}

func calculateActualBotShotSide(
	botInfo PlayerInfo,
	shotType ShotType,
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
		return botCalculateActualAttackShotSide(botInfo, botShotSide, playerInfo, playerShotSide)

	case ShotTypeDefend:
		return botCalculateActualDefendShotSide(botInfo, botShotSide, playerInfo, playerShotSide)
	}

	return botShotSide
}

func botCalculateActualAttackShotSide(
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

func botCalculateActualDefendShotSide(
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
	for _, s := range possibleBotShotSides {
		if s != side {
			return s
		}
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

func generateBotSide() ShotSide {
	//nolint:gosec
	return possibleBotShotSides[rand.Int()%len(possibleBotShotSides)]
}

// playersShots returns in shot by player id for modification, shot another player for ready only.
func (g *Game) playersShots(id player.ID) (*Shot, Shot) {
	cRound := g.currentRoundPtr()

	if cRound.Attack.PlayerID == id {
		return &cRound.Attack, cRound.Defend
	}

	return &cRound.Defend, cRound.Attack
}
