package game

import (
	"errors"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	ShotsInitial = 5

	shotsToCompleteRound = 2
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type Game struct {
	ID              ID
	CreatedAt       time.Time
	Status          GameStatus
	StatusChagnedAt time.Time
	State           State
	StateVersion    int
}

type Player struct {
	ID             player.ID
	DisplayName    string
	ShotsAvailable int
	GoalsScored    int
}

type RoundShot struct {
	PlayerID player.ID
	Side     ShotSide
}

type Round struct {
	Defend RoundShot
	Attack RoundShot
	IsGoal bool
}

type State struct {
	Players []Player
	Rounds  []Round
}

func (g *Game) CurrentRound() int {
	return len(g.State.Rounds)
}

func (g *Game) IsFinished() bool {
	return g.Status == GameStatusFinished
}

func (g *Game) CompleteRound(shots []Shot) error {
	if len(shots) != shotsToCompleteRound {
		return errors.New("invalid shot count")
	}

	round := makeRoundByShots(shots)

	g.State.Rounds = append(g.State.Rounds, round)

	g.updateAttackerShots(round)

	g.updateGameStatus()

	return nil
}

func makeRoundByShots(shots []Shot) Round {
	var round Round

	for _, shot := range shots {
		roundShot := RoundShot{
			PlayerID: shot.PlayerID,
			Side:     shot.Side,
		}

		switch shot.Type {
		case ShotTypeAttack:
			round.Attack = roundShot

		case ShotTypeDefend:
			round.Defend = roundShot
		}
	}

	round.IsGoal = round.Attack.Side != round.Defend.Side

	return round
}

func (g *Game) updateAttackerShots(round Round) {
	for plrIdx, plr := range g.State.Players {
		if plr.ID != round.Attack.PlayerID {
			continue
		}

		g.State.Players[plrIdx].ShotsAvailable--

		if round.IsGoal {
			g.State.Players[plrIdx].GoalsScored++
		}
	}
}

func (g *Game) updateGameStatus() {
	for _, plr := range g.State.Players {
		if plr.ShotsAvailable > 0 {
			return
		}
	}

	g.Status = GameStatusFinished
}
