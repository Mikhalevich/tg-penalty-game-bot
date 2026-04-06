package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type PlayerProvider interface {
	GetPlayerByChatID(ctx context.Context, chatID msginfo.ChatID) (player.Player, error)
}

type Welcome interface {
	Welcome(ctx context.Context, chatID msginfo.ChatID) error
}

type ChangeName interface {
	SetChangeDisplayNameTrigger(ctx context.Context, chatID msginfo.ChatID, fullName, userName string) error
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, msgID msginfo.MessageID, displayName string) error
}

type FindGame interface {
	FindGame(ctx context.Context, chatID msginfo.ChatID) error
	StopFind(ctx context.Context, chatID msginfo.ChatID) error
}

type StartGame interface {
	StartGameWithBot(ctx context.Context, chatID msginfo.ChatID) error
	StartGameByLink(ctx context.Context, chatID msginfo.ChatID) error
	JoinGameByLink(ctx context.Context, chatID msginfo.ChatID, gameID game.ID) error
}

type LeaveGame interface {
	LeaveGame(ctx context.Context, chatID msginfo.ChatID) error
}

type Leaderboard interface {
	Page(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID, pageNumber int) error
}

type GameShot interface {
	Shot(
		ctx context.Context,
		chatID msginfo.ChatID,
		msgID msginfo.MessageID,
		gameID game.ID,
		round int,
		side game.ShotSide,
	) error
}

type TGHandler struct {
	buttonProvider ButtonProvider
	playerProvider PlayerProvider
	welcome        Welcome
	changeName     ChangeName
	findGame       FindGame
	startGame      StartGame
	gameShot       GameShot
	leaveGame      LeaveGame
	leaderboard    Leaderboard
}

func New(
	buttonProvider ButtonProvider,
	playerProvider PlayerProvider,
	welcome Welcome,
	changeName ChangeName,
	findGame FindGame,
	startGame StartGame,
	gameShot GameShot,
	leaveGame LeaveGame,
	leaderboard Leaderboard,
) *TGHandler {
	return &TGHandler{
		buttonProvider: buttonProvider,
		playerProvider: playerProvider,
		welcome:        welcome,
		changeName:     changeName,
		findGame:       findGame,
		startGame:      startGame,
		gameShot:       gameShot,
		leaveGame:      leaveGame,
		leaderboard:    leaderboard,
	}
}
