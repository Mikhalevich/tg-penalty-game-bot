package game

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

// CreateGame create a new game(in progress status), start first round and return created game.
func CreateGame(
	ctx context.Context,
	gameType GameType,
	player1 player.Player,
	player2 player.Player,
	createdAt time.Time,
) Game {
	gameID := IDFromString(uuid.NewString())
	createdGame := Game{
		ID:        gameID,
		CreatedAt: createdAt,
		Type:      gameType,
		Status:    GameStatusInProgress,
		State: State{
			Player1: CreateGamePlayerFromPlayer(player1),
			Player2: CreateGamePlayerFromPlayer(player2),
			Rounds:  makeRounds(gameID, ShotsInitial, player1.ID, player2.ID),
		},
		StateUpdatedAt: createdAt,
	}

	createdGame.StartFirstRound(createdAt)

	return createdGame
}

// CreatePendingGame create a new game in pending status, no round is starting and return created game.
func CreatePendingGame(
	ctx context.Context,
	plr player.Player,
	createdAt time.Time,
) Game {
	return Game{
		ID:        IDFromString(uuid.NewString()),
		CreatedAt: createdAt,
		Type:      GameTypeFriendly,
		Status:    GameStatusPending,
		State: State{
			Player1: CreateGamePlayerFromPlayer(plr),
		},
		StateUpdatedAt: createdAt,
	}
}

func CreateGamePlayerFromPlayer(plr player.Player) Player {
	return Player{
		ID:          plr.ID,
		ChatID:      plr.ChatID,
		DisplayName: plr.DisplayName,
	}
}

func makeRounds(
	gameID ID,
	shotsCount int,
	playerID1 player.ID,
	playerID2 player.ID,
) []Round {
	var (
		roundsCount = shotsCount * 2
		rounds      = make([]Round, 0, roundsCount)
		attackerID  = playerID1
		defenderID  = playerID2
	)
	for roundNumber := range roundsCount {
		rounds = append(rounds, Round{
			Attack: makePendingShot(roundNumber, gameID, attackerID, ShotTypeAttack),
			Defend: makePendingShot(roundNumber, gameID, defenderID, ShotTypeDefend),
			Result: RoundResultNotCompleted,
		})

		attackerID, defenderID = defenderID, attackerID
	}

	return rounds
}

func makePendingShot(
	roundNumber int,
	gameID ID,
	playerID player.ID,
	shotType ShotType,
) Shot {
	return Shot{
		GameID:   gameID,
		PlayerID: playerID,
		Round:    roundNumber,
		Type:     shotType,
		Side:     ShotSideNoShot,
	}
}
