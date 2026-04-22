package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
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
	ChangeDisplayName(ctx context.Context, chatID msginfo.ChatID, displayName string) error
	Cancel(ctx context.Context, chatID msginfo.ChatID) error
}

type FindGame interface {
	FindGame(ctx context.Context, chatID msginfo.ChatID) error
	StopFind(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
}

type StartGame interface {
	StartGameWithBot(ctx context.Context, chatID msginfo.ChatID) error
	StartGameByLink(ctx context.Context, chatID msginfo.ChatID) error
	JoinGameByLink(ctx context.Context, chatID msginfo.ChatID, gameID game.ID) error
}

type LeaveGame interface {
	LeaveGame(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error
}

type Leaderboard interface {
	Page(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID, pageNumber int) error
}

type ShotStats interface {
	ViewStatsOnStartGameMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		playerID player.ID,
		playerDisplayName string,
	) error
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

type Notifier interface {
	ParseError(ctx context.Context, chatID msginfo.ChatID, err error) error
}

type cbHandler func(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error

type TGHandler struct {
	cbHanlers      map[button.Operation]cbHandler
	buttonProvider ButtonProvider
	playerProvider PlayerProvider
	welcome        Welcome
	changeName     ChangeName
	findGame       FindGame
	startGame      StartGame
	gameShot       GameShot
	leaveGame      LeaveGame
	leaderboard    Leaderboard
	shotStats      ShotStats
	notifier       Notifier
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
	shotStats ShotStats,
	notifier Notifier,
) *TGHandler {
	tgh := &TGHandler{
		buttonProvider: buttonProvider,
		playerProvider: playerProvider,
		welcome:        welcome,
		changeName:     changeName,
		findGame:       findGame,
		startGame:      startGame,
		gameShot:       gameShot,
		leaveGame:      leaveGame,
		leaderboard:    leaderboard,
		shotStats:      shotStats,
		notifier:       notifier,
	}

	tgh.registerCBHandlers()

	return tgh
}

func (t *TGHandler) registerCBHandlers() {
	t.cbHanlers = map[button.Operation]cbHandler{
		button.OperationChangeName:        t.cbChangeName,
		button.OperationChangeNameTrigger: t.cbChangeNameTrigger,
		button.OperationChangeNameCancel:  t.cbChangeNameCancel,
		button.OperationShotSide:          t.cbShotSide,
		button.OperationLeaveGame:         t.cbLeaveGame,
		button.OperationShotStats:         t.cbShotStatsOnStartGameMessage,
		button.OperationStopSearchGame:    t.cbStopFindButton,
		button.OperationLeaderboardPage:   t.cbLeaderboardPage,
	}
}
