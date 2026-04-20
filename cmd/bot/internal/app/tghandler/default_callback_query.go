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

	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Button expired")

			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	hndlr, ok := t.cbHanlers[btn.Operation]
	if !ok {
		return fmt.Errorf("invalid operation %v", btn.Operation)
	}

	if err := hndlr(ctx, msg, btn); err != nil {
		return fmt.Errorf("process cb handler: %w", err)
	}

	return nil
}

func (t *TGHandler) cbChangeName(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[button.ChangeNamePayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := t.changeDisplayName(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.DisplayName,
	); err != nil {
		return fmt.Errorf("change display name: %w", err)
	}

	return nil
}

func (t *TGHandler) cbChangeNameTrigger(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.changeName.SetChangeDisplayNameTrigger(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msg.User.FullName(),
		msg.User.Username,
	); err != nil {
		return fmt.Errorf("set change name trigger: %w", err)
	}

	return nil
}

func (t *TGHandler) cbChangeNameCancel(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.changeName.Cancel(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("change name cancel: %w", err)
	}

	return nil
}

func (t *TGHandler) cbShotSide(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[game.ShotSidePayload](*btn)
	if err != nil {
		return fmt.Errorf("get payload: %w", err)
	}

	if err := t.gameShot.Shot(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.GameID,
		payload.Round,
		payload.Side,
	); err != nil {
		return fmt.Errorf("shot: %w", err)
	}

	return nil
}

func (t *TGHandler) cbLeaveGame(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.leaveGame.LeaveGame(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
	); err != nil {
		return fmt.Errorf("leave game: %w", err)
	}

	return nil
}

func (t *TGHandler) cbShotStatsOnStartGameMessage(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[game.ShotStatsPayload](*btn)
	if err != nil {
		return fmt.Errorf("stats payload: %w", err)
	}

	if err := t.shotStats.ViewStatsOnStartGameMessage(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.PlayerID,
		payload.DisplayName,
	); err != nil {
		return fmt.Errorf("view stats on start game msg: %w", err)
	}

	return nil
}

func (t *TGHandler) cbStopFindButton(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.findGame.StopFind(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
	); err != nil {
		return fmt.Errorf("stop find: %w", err)
	}

	return nil
}

func (t *TGHandler) cbLeaderboardPage(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[button.LeaderboardPagePayload](*btn)
	if err != nil {
		return fmt.Errorf("leaderboard page payload: %w", err)
	}

	if err := t.leaderboard.Page(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.PageNumber,
	); err != nil {
		return fmt.Errorf("leaderboard player: %w", err)
	}

	return nil
}
