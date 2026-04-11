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
	ID             ID
	CreatedAt      time.Time
	Type           GameType
	Status         GameStatus
	State          State
	StateVersion   int
	StateUpdatedAt time.Time
}

type Player struct {
	ID          player.ID
	ChatID      msginfo.ChatID
	DisplayName string
	GoalsScored int
}

func (p Player) IsBot() bool {
	return p.ID == 0
}

func (g *Game) IsGameWithBot() bool {
	return g.State.Player1.IsBot() || g.State.Player2.IsBot()
}

func (g *Game) IsFinished() bool {
	return g.Status == GameStatusCompleted
}

func (g *Game) IsInProgress() bool {
	return g.Status == GameStatusInProgress
}

func (g *Game) IsRatingGame() bool {
	return g.Type == GameTypeRating
}

func (g *Game) Cancel(canceledAt time.Time) {
	g.Status = GameStatusCanceled
	g.StateUpdatedAt = canceledAt
}

func (g *Game) PlayerIDs() []player.ID {
	//nolint:mnd
	ids := make([]player.ID, 0, 2)

	if !g.State.Player1.IsBot() {
		ids = append(ids, g.State.Player1.ID)
	}

	if !g.State.Player2.IsBot() {
		ids = append(ids, g.State.Player2.ID)
	}

	return ids
}

func (g *Game) JoinPlayerAndStartGame(plr Player, joinedAt time.Time) error {
	if g.Status != GameStatusPending {
		return perror.InvalidGameState()
	}

	if g.State.Player1.ID == plr.ID {
		return perror.InvalidPlayer()
	}

	g.State.Player2 = plr
	g.State.Rounds = makeRounds(g.ID, ShotsInitial, g.State.Player1.ID, g.State.Player2.ID)

	g.StartFirstRound(joinedAt)

	g.Status = GameStatusInProgress
	g.StateUpdatedAt = joinedAt

	return nil
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

	g.StateUpdatedAt = shot.CompletedAt

	return nil
}

// TryToCompleteRound complete round if both players make shots
// and updates attacker shots and scores
// returns true if round was completed.
func (g *Game) TryToCompleteRound() bool {
	cRound := g.currentRoundPtr()

	if (cRound.Attack.Side == ShotSideNoShot) ||
		(cRound.Defend.Side == ShotSideNoShot) {
		return false
	}

	updateRoundResults(cRound)

	g.updateAttackerGoals(cRound.Attack.PlayerID, cRound.Result == RoundResultGoal)

	return true
}

// StartFirstRound just update created_at shot times.
func (g *Game) StartFirstRound(startedAt time.Time) {
	cRound := g.currentRoundPtr()
	cRound.Attack.CreatedAt = startedAt
	cRound.Defend.CreatedAt = startedAt
}

// StartNextRound starts next round or finish the game.
func (g *Game) StartNextRound(startedAt time.Time) error {
	if !g.State.CurrentRound().IsCompleted() {
		return perror.RoundNotCompleted()
	}

	g.State.CurrentRoundIdx++
	g.StateUpdatedAt = startedAt

	if g.State.CurrentRoundIdx >= len(g.State.Rounds) {
		g.Status = GameStatusCompleted

		return nil
	}

	cRound := g.currentRoundPtr()
	cRound.Attack.CreatedAt = startedAt
	cRound.Defend.CreatedAt = startedAt

	return nil
}

func (g *Game) MakeMissShots(shotAt time.Time) {
	round := g.currentRoundPtr()

	updateIfNoShot(&round.Attack, ShotSideMiss, shotAt)
	updateIfNoShot(&round.Defend, ShotSideMiss, shotAt)

	g.StateUpdatedAt = shotAt
}

func (g *Game) MissForRestShotsAndCompleteGame(playerID player.ID, completedAt time.Time) {
	for i := g.State.CurrentRoundIdx; i < len(g.State.Rounds); i++ {
		round := &g.State.Rounds[i]
		missShotForPlayerOrMiddleOtherwise(&round.Attack, playerID, completedAt)
		missShotForPlayerOrMiddleOtherwise(&round.Defend, playerID, completedAt)

		updateRoundResults(round)

		g.updateAttackerGoals(round.Attack.PlayerID, round.Result == RoundResultGoal)
	}

	g.Status = GameStatusCompleted
	g.StateUpdatedAt = completedAt
}

func (g *Game) currentRoundPtr() *Round {
	return &g.State.Rounds[g.State.CurrentRoundIdx]
}

func updateShot(shot *Shot, side ShotSide, shotAt time.Time) {
	shot.Side = side
	shot.CompletedAt = shotAt
}

func updateIfNoShot(shot *Shot, side ShotSide, shotAt time.Time) {
	if shot.Side != ShotSideNoShot || shot.PlayerID == 0 {
		return
	}

	updateShot(shot, side, shotAt)
}

func missShotForPlayerOrMiddleOtherwise(
	shot *Shot,
	missShotPlayerID player.ID,
	completedAt time.Time,
) {
	if shot.PlayerID == missShotPlayerID {
		updateShot(shot, ShotSideMiss, completedAt)

		return
	}

	updateIfNoShot(shot, ShotSideMiddle, completedAt)
}

func (g *Game) shotByPlayerID(id player.ID) *Shot {
	cRound := g.currentRoundPtr()

	if cRound.Attack.PlayerID == id {
		return &cRound.Attack
	}

	return &cRound.Defend
}

func (g *Game) updateAttackerGoals(attackerID player.ID, isGoal bool) {
	var attackerPlayer *Player

	if attackerID == g.State.Player1.ID {
		attackerPlayer = &g.State.Player1
	} else {
		attackerPlayer = &g.State.Player2
	}

	if isGoal {
		attackerPlayer.GoalsScored++
	}
}

func updateRoundResults(round *Round) {
	switch {
	case round.Attack.Side == ShotSideMiss:
		round.Result = RoundResultMiss

	case round.Attack.Side != round.Defend.Side:
		round.Result = RoundResultGoal

	default:
		round.Result = RoundResultSave
	}
}
