package tournament

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
)

type Team struct {
	Player Player
	Stat   GameStat
}

type GameStat struct {
	Played        int
	Points        int
	GoalsScored   int
	GoalsConceded int
}

type Match struct {
	ID          game.ID
	IsCompleted bool
	Home        MatchTeamInfo
	Away        MatchTeamInfo
}

type MatchTeamInfo struct {
	Idx         int
	GoalsScored int
}

type Round []Match

type State struct {
	Teams        []Team
	Rounds       []Round
	CurrentRound int
}
