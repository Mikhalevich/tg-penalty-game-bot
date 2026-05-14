package game

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	percent100 = 100
	percent60  = 60
	percent25  = 25
	percent30  = 30
	percent20  = 20
	percent10  = 10
	percent1   = 1
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
		DefendMissPercent: percent1,
	}
}

func createPlayerInfoAgainstBot(difficulty BotDifficulty) PlayerInfo {
	switch difficulty {
	case BotDifficultyEasy:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent1,
			ForceGoalPercent:  percent30,
			ForceSavePercent:  percent30,
		}

	case BotDifficultyNormal:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent1,
			ForceGoalPercent:  percent10,
			ForceSavePercent:  percent10,
		}

	case BotDifficultyHard, BotDifficultyInsane:
		return PlayerInfo{
			AttackMissPercent: percent10,
			DefendMissPercent: percent1,
		}
	}

	return PlayerInfo{
		AttackMissPercent: percent10,
		DefendMissPercent: percent1,
	}
}

//nolint:funlen
func createBot(difficulty BotDifficulty) Player {
	switch difficulty {
	case BotDifficultyEasy:
		return Player{
			ID:          0,
			DisplayName: "bot easy",
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent20,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}

	case BotDifficultyNormal:
		return Player{
			ID:          0,
			DisplayName: "bot normal",
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent20,
				DefendMissPercent: percent1,
				ForceGoalPercent:  0,
				ForceSavePercent:  0,
			},
		}

	case BotDifficultyHard:
		return Player{
			ID:          0,
			DisplayName: "bot hard",
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent1,
				DefendMissPercent: percent1,
				ForceGoalPercent:  percent30,
				ForceSavePercent:  percent25,
			},
		}

	case BotDifficultyInsane:
		return Player{
			ID:          0,
			DisplayName: "bot insane",
			Info: PlayerInfo{
				IsBot:             true,
				AttackMissPercent: percent1,
				DefendMissPercent: percent1,
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
			DefendMissPercent: percent1,
			ForceGoalPercent:  0,
			ForceSavePercent:  0,
		},
	}
}
