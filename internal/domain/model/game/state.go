package game

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type RoundResult int

const (
	RoundResultNotCompleted RoundResult = iota
	RoundResultGoal
	RoundResultSave
	RoundResultMiss
)

type Round struct {
	Defend Shot
	Attack Shot
	Result RoundResult
}

func (r Round) IsCompleted() bool {
	return r.Result != RoundResultNotCompleted
}

func (r Round) IsGoal() bool {
	return r.Result == RoundResultGoal
}

type State struct {
	Player1         Player
	Player2         Player
	Rounds          []Round
	CurrentRoundIdx int
}

func (s State) CurrentRound() Round {
	return s.Rounds[s.CurrentRoundIdx]
}

func (s State) CurrentRoundNumber() int {
	return s.CurrentRoundIdx
}

func (s State) LivePlayers() []Player {
	if s.Player1.IsBot() && s.Player2.IsBot() {
		return nil
	}

	if s.Player1.IsBot() {
		return []Player{s.Player2}
	}

	if s.Player2.IsBot() {
		return []Player{s.Player1}
	}

	return []Player{s.Player1, s.Player2}
}

func (s State) PlayerByID(playerID player.ID) Player {
	if s.Player1.ID == playerID {
		return s.Player1
	}

	return s.Player2
}
