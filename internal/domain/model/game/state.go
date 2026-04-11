package game

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	shotsPerRound = 2
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

func (s State) CompletedShots() []Shot {
	shots := make([]Shot, 0, len(s.Rounds)*shotsPerRound)

	for _, round := range s.Rounds {
		if !round.IsCompleted() {
			continue
		}

		shots = append(shots, round.Attack, round.Defend)
	}

	return shots
}
