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

func (s State) isRoundsCompleted() bool {
	return s.CurrentRound >= len(s.Rounds)
}

func (s State) isCurrentRoundFinished() bool {
	if s.isRoundsCompleted() {
		return true
	}

	for _, match := range s.Rounds[s.CurrentRound] {
		if !match.IsCompleted {
			return false
		}
	}

	return true
}
