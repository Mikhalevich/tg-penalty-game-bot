package tournament

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

const (
	MinPlayers = 3
	MaxPlayers = 20
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(id string) ID {
	return ID(id)
}

type TournamentStatus string

const (
	TournamentStatusPending    TournamentStatus = "pending"
	TournamentStatusInProgress TournamentStatus = "in_progress"
	TournamentStatusCompleted  TournamentStatus = "completed"
	TournamentStatusCanceled   TournamentStatus = "canceled"
)

func (ts TournamentStatus) String() string {
	return string(ts)
}

type Tournament struct {
	ID             ID
	CreatedAt      time.Time
	Status         TournamentStatus
	State          State
	StateVersion   int
	StateUpdatedAt time.Time
}

func (t *Tournament) Join(plr Player) error {
	if err := t.ensureStatus(TournamentStatusPending); err != nil {
		return fmt.Errorf("ensure status: %w", err)
	}

	if len(t.State.Teams) >= MaxPlayers {
		return perror.InvalidState("max players reached")
	}

	t.State.Teams = append(t.State.Teams, Team{
		Player: plr,
	})

	return nil
}

func (t *Tournament) Start(startedAt time.Time) error {
	if err := t.ensureStatus(TournamentStatusPending); err != nil {
		return fmt.Errorf("ensure status: %w", err)
	}

	if len(t.State.Teams) < MinPlayers {
		return perror.InvalidState("not enough players")
	}

	t.State.Rounds = RoundRobinSchedule(len(t.State.Teams))
	t.Status = TournamentStatusInProgress
	t.StateUpdatedAt = startedAt

	return nil
}

// NextRound try to start next round
// returns new round games, completion flag and error.
func (t *Tournament) NextRound(startedAt time.Time) ([]game.Game, bool, error) {
	if !t.State.isCurrentRoundFinished() {
		return nil, false, perror.InvalidState("current round is not finished")
	}

	t.State.CurrentRound++

	if t.State.isRoundsCompleted() {
		t.Status = TournamentStatusCompleted
		t.StateUpdatedAt = startedAt

		return nil, true, nil
	}

	games := t.makeRoundGames(startedAt)

	return games, false, nil
}

func (t *Tournament) makeRoundGames(startedAt time.Time) []game.Game {
	var (
		games               = make([]game.Game, 0, len(t.State.Rounds[t.State.CurrentRound]))
		currentRoundMatches = t.State.Rounds[t.State.CurrentRound]
	)

	for matchIdx, match := range currentRoundMatches {
		var (
			homePlayer = t.State.Teams[match.Home.Idx]
			awayPlayer = t.State.Teams[match.Away.Idx]
		)

		newGame := makeGame(homePlayer.Player, awayPlayer.Player, startedAt)
		currentRoundMatches[matchIdx].ID = newGame.ID
		games = append(games, newGame)
	}

	return games
}

func makeGame(homePlayer, awayPlayer Player, startedAt time.Time) game.Game {
	switch {
	case homePlayer.IsBot && awayPlayer.IsBot:
		return game.CreateGameBotToBot(
			game.GameTypeTournament,
			homePlayer.BotDifficulty,
			awayPlayer.BotDifficulty,
			startedAt,
		)

	case homePlayer.IsBot:
		return game.CreateGameWithBot(
			game.GameTypeTournament,
			awayPlayer.toDomPlayer(),
			homePlayer.BotDifficulty,
			startedAt,
		)

	case awayPlayer.IsBot:
		return game.CreateGameWithBot(
			game.GameTypeTournament,
			homePlayer.toDomPlayer(),
			awayPlayer.BotDifficulty,
			startedAt,
		)

	default:
		return game.CreateGameWithPlayer(
			game.GameTypeTournament,
			homePlayer.toDomPlayer(),
			awayPlayer.toDomPlayer(),
			startedAt,
		)
	}
}

func (t *Tournament) ensureStatus(s TournamentStatus) error {
	if t.Status == s {
		return nil
	}

	switch t.Status {
	case TournamentStatusPending:
		return perror.InvalidState("tournament is in pending state")

	case TournamentStatusInProgress:
		return perror.InvalidState("tournament already in progress")

	case TournamentStatusCompleted:
		return perror.InvalidState("tournament is completed")

	case TournamentStatusCanceled:
		return perror.InvalidState("tournament is canceled")
	}

	return nil
}
