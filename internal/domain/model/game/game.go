package game

import (
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	ShotsInitial = 5
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
	ChatID         msginfo.ChatID
	DisplayName    string
	ShotsAvailable int
	GoalsScored    int
}

type Round struct {
	Defend      Shot
	Attack      Shot
	IsGoal      bool
	IsCompleted bool
}

type State struct {
	Player1 Player
	Player2 Player
	Rounds  []Round
}

func (g *Game) CurrentRoundNumber() int {
	return len(g.State.Rounds)
}

func (g *Game) CurrentRound() *Round {
	if len(g.State.Rounds) == 0 {
		return nil
	}

	return &g.State.Rounds[len(g.State.Rounds)-1]
}

func (g *Game) IsFinished() bool {
	return g.Status == GameStatusFinished
}

func (g *Game) PlayerShot(shot Shot) error {
	inGameShot := g.shotByPlayerID(shot.PlayerID)

	if inGameShot.Side != ShotSideNoShot {
		return perror.AlreadyExists("shot already exist")
	}

	if inGameShot.PlayerID != shot.PlayerID {
		return perror.InvalidPlayer()
	}

	if inGameShot.Round != shot.Round {
		return perror.InvalidRound()
	}

	inGameShot.Side = shot.Side
	inGameShot.CompletedAt = shot.CompletedAt

	return nil
}

// TryToCompleteRound complete round if both players make shots
// and updates attacker shots and scores
// returns true if round was completed.
func (g *Game) TryToCompleteRound() bool {
	cRound := g.CurrentRound()

	if (cRound.Attack.Side == ShotSideNoShot) ||
		(cRound.Defend.Side == ShotSideNoShot) {
		return false
	}

	cRound.IsGoal = cRound.Attack.Side != cRound.Defend.Side

	g.updateAttackerShots(cRound.Attack.PlayerID, cRound.IsGoal)

	cRound.IsCompleted = true

	return true
}

// StartNextRound starts next round or finish the game.
func (g *Game) StartNextRound(now time.Time) error {
	cRound := g.CurrentRound()
	if cRound != nil && !cRound.IsCompleted {
		return perror.RoundNotCompleted()
	}

	var (
		player1 = g.State.Player1
		player2 = g.State.Player2
	)

	if (player1.ShotsAvailable == 0) && (player2.ShotsAvailable == 0) {
		g.Status = GameStatusFinished

		return nil
	}

	attacker, defender := nextAttackerDefenderInOrder(player1, player2)

	g.State.Rounds = append(g.State.Rounds, Round{
		Attack: g.makePendingShot(attacker.ID, ShotTypeAttack, now),
		Defend: g.makePendingShot(defender.ID, ShotTypeDefend, now),
	})

	return nil
}

func (g *Game) shotByPlayerID(id player.ID) *Shot {
	cRound := g.CurrentRound()

	if cRound.Attack.PlayerID == id {
		return &cRound.Attack
	}

	return &cRound.Defend
}

func (g *Game) makePendingShot(playerID player.ID, shotType ShotType, now time.Time) Shot {
	return Shot{
		GameID:    g.ID,
		PlayerID:  playerID,
		Round:     g.CurrentRoundNumber() + 1,
		Type:      shotType,
		CreatedAt: now,
		Side:      ShotSideNoShot,
	}
}

// nextAttackerDefenderInOrder returns players in order attacker => defender.
func nextAttackerDefenderInOrder(player1, player2 Player) (Player, Player) {
	if player1.ShotsAvailable >= player2.ShotsAvailable {
		return player1, player2
	}

	return player2, player1
}

func (g *Game) updateAttackerShots(attackerID player.ID, isGoal bool) {
	var attackerPlayer *Player

	if attackerID == g.State.Player1.ID {
		attackerPlayer = &g.State.Player1
	} else {
		attackerPlayer = &g.State.Player2
	}

	attackerPlayer.ShotsAvailable--

	if isGoal {
		attackerPlayer.GoalsScored++
	}
}
