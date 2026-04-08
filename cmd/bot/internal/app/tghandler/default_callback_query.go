package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Data == "" {
		return nil
	}

	var (
		chatID = msginfo.ChatIDFromInt64(msg.ChatID)
		msgID  = msginfo.MessageIDFromInt(msg.MessageID)
	)

	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Button expired")

			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	switch btn.Operation {
	case button.OperationChangeName:
		return t.processChangeNameButton(ctx, chatID, msgID, btn)

	case button.OperationShotSide:
		return t.processShotSideButton(ctx, chatID, msgID, btn)

	case button.OperationLeaveGame:
		return t.processLeaveGameButton(ctx, chatID)

	case button.OperationStopSearchGame:
		return t.processStopFindButton(ctx, chatID)

	case button.OperationLeaderboardPage:
		return t.processLeaderboardPage(ctx, chatID, msginfo.MessageIDFromInt(msg.MessageID), btn)
	}

	return nil
}

func (t *TGHandler) processChangeNameButton(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[button.ChangeNamePayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := t.changeDisplayName(ctx, chatID, msgID, payload.DisplayName); err != nil {
		return fmt.Errorf("change display name: %w", err)
	}

	return nil
}

func (t *TGHandler) processShotSideButton(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[game.ShotSidePayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := t.gameShot.Shot(ctx, chatID, msgID, payload.GameID, payload.Round, payload.Side); err != nil {
		return fmt.Errorf("shot: %w", err)
	}

	return nil
}

func (t *TGHandler) processLeaveGameButton(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := t.leaveGame.LeaveGame(ctx, chatID); err != nil {
		return fmt.Errorf("leave game: %w", err)
	}

	return nil
}

func (t *TGHandler) processStopFindButton(
	ctx context.Context,
	chatID msginfo.ChatID,
) error {
	if err := t.findGame.StopFind(ctx, chatID); err != nil {
		return fmt.Errorf("stop find: %w", err)
	}

	return nil
}

func (t *TGHandler) processLeaderboardPage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[button.LeaderboardPagePayload](*btn)
	if err != nil {
		return fmt.Errorf("leaderboard page payload: %w", err)
	}

	if err := t.leaderboard.Page(ctx, chatID, messageID, payload.PageNumber); err != nil {
		return fmt.Errorf("leaderboard player: %w", err)
	}

	return nil
}
