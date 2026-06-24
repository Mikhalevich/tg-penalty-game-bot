package game

import (
	"time"
)

// CreateGameBotToBot create a bot to bot game(in progress status).
func CreateGameBotToBot(
	gameType GameType,
	homeDiff BotDifficulty,
	awayDiff BotDifficulty,
	createdAt time.Time,
) Game {
	var (
		gameID  = GenerateID()
		homeBot = createBot(homeDiff)
		awayBot = createBot(awayDiff)
	)

	adjustBotParams(&homeBot, &awayBot)

	createdGame := Game{
		ID:        gameID,
		CreatedAt: createdAt,
		Type:      gameType,
		Status:    GameStatusInProgress,
		State: State{
			Player1: homeBot,
			Player2: awayBot,
			Rounds:  makeRounds(gameID, ShotsInitial, homeBot.ID, awayBot.ID),
		},
		StateUpdatedAt: createdAt,
	}

	createdGame.makeGameBetweenBots(createdAt)

	return createdGame
}

func (g *Game) makeGameBetweenBots(startedAt time.Time) {
	for range g.State.Rounds {
		cRound := g.currentRoundPtr()
		cRound.Attack.CreatedAt = startedAt
		cRound.Defend.CreatedAt = startedAt

		cRound.Attack.ExpectedSide = generateExpectedShotSide()
		cRound.Attack.ActualSide = cRound.Attack.ExpectedSide

		cRound.Defend.ExpectedSide = generateExpectedShotSide()
		cRound.Defend.ActualSide = cRound.Defend.ExpectedSide

		updateRoundResults(cRound)

		cRound.Attack.CompletedAt = startedAt
		cRound.Defend.CompletedAt = startedAt

		g.State.CurrentRoundIdx++
	}

	g.StateUpdatedAt = startedAt
	g.Status = GameStatusCompleted
}

func adjustBotParams(homeBot, awayBot *Player) {
	var (
		forceGoalDiff = homeBot.Info.ForceGoalPercent - awayBot.Info.ForceGoalPercent
		forceSaveDiff = homeBot.Info.ForceSavePercent - awayBot.Info.ForceSavePercent
	)

	homeBot.Info.ForceGoalPercent = positive(forceGoalDiff)
	awayBot.Info.ForceGoalPercent = positive(-forceGoalDiff)

	homeBot.Info.ForceSavePercent = positive(forceSaveDiff)
	awayBot.Info.ForceSavePercent = positive(-forceSaveDiff)
}

func positive(value int) int {
	if value > 0 {
		return value
	}

	return 0
}
