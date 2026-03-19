package game

import (
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type Round struct {
	Defend      Shot
	Attack      Shot
	IsGoal      bool
	IsCompleted bool
}

type Rounds []Round

func (r Rounds) Last() Round {
	if len(r) == 0 {
		return Round{}
	}

	return r[len(r)-1]
}

func (r Rounds) Len() int {
	return len(r)
}

type State struct {
	Player1 Player
	Player2 Player
	Rounds  Rounds
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
