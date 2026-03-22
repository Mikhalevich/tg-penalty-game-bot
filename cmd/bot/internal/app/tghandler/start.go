package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tghandler/internal/cmdargs"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (t *TGHandler) Start(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	operation, args := cmdargs.Parse(msg.Args)

	switch operation {
	case cmdargs.OperationJoin:
		if err := t.startGame.JoinGameByLink(
			ctx,
			msginfo.ChatIDFromInt64(msg.ChatID),
			game.IDFromString(args),
		); err != nil {
			return fmt.Errorf("join game by link: %w", err)
		}

	case cmdargs.OperationNope:
		if err := t.welcome.Welcome(ctx, msginfo.ChatIDFromInt64(msg.ChatID)); err != nil {
			return fmt.Errorf("welcome: %w", err)
		}
	}

	return nil
}
