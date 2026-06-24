package game

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	percent100 = 100
	percent60  = 60
	percent40  = 40
	percent30  = 30
	percent25  = 25
	percent20  = 20
	percent10  = 10
	percent2   = 2
)

const (
	botEasy   = "bot easy"
	botNormal = "bot normal"
	botHard   = "bot hard"
	botInsane = "bot insane"
)

type BotDifficulty int

const (
	BotDifficultyEasy BotDifficulty = iota + 1
	BotDifficultyNormal
	BotDifficultyHard
	BotDifficultyInsane
)

type Player struct {
	ID          player.ID
	ChatID      msginfo.ChatID
	DisplayName string
	GoalsScored int
	Info        PlayerInfo
}

type PlayerInfo struct {
	IsBot             bool
	AttackMissPercent int
	DefendMissPercent int
	ForceGoalPercent  int
	ForceSavePercent  int
}

func (p Player) IsBot() bool {
	return p.Info.IsBot
}

func createPlayerInfo() PlayerInfo {
	return PlayerInfo{
		AttackMissPercent: percent10,
		DefendMissPercent: percent2,
	}
}

func createPlayerInfoAgainstBot(difficulty BotDifficulty) PlayerInfo {
	switch difficulty {
	case BotDifficultyEasy:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent2,
			ForceGoalPercent:  percent30,
			ForceSavePercent:  percent40,
		}

	case BotDifficultyNormal:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent2,
			ForceGoalPercent:  0,
			ForceSavePercent:  percent20,
		}

	case BotDifficultyHard, BotDifficultyInsane:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent2,
		}
	}

	return PlayerInfo{
		AttackMissPercent: percent10,
		DefendMissPercent: percent2,
	}
}

//nolint:funlen
func createBot(difficulty BotDifficulty) Player {
	switch difficulty {
	case BotDifficultyEasy:
		return Player{
			ID:          0,
			DisplayName: botEasy,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent10,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}

	case BotDifficultyNormal:
		return Player{
			ID:          0,
			DisplayName: botNormal,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent2,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}

	case BotDifficultyHard:
		return Player{
			ID:          0,
			DisplayName: botHard,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  percent30,
				ForceSavePercent:  percent25,
			},
		}

	case BotDifficultyInsane:
		return Player{
			ID:          0,
			DisplayName: botInsane,
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent2,
				DefendMissPercent: percent2,
				ForceGoalPercent:  percent60,
				ForceSavePercent:  percent60,
			},
		}
	}

	return Player{
		ID:          0,
		DisplayName: "bot",
		Info: PlayerInfo{
			IsBot:             true,
			AttackMissPercent: percent10,
			DefendMissPercent: percent2,
			ForceGoalPercent:  0,
			ForceSavePercent:  0,
		},
	}
}
