package viewstats

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (vs *ViewStats) ViewStatsOnStartGameMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
) error {
	currentPlayer, err := vs.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player by chat_id: %w", err)
	}

	if currentPlayer.GameStatus != player.GameStatusInGame {
		if err := vs.messageDeleter.DeleteMessage(ctx, chatID, messageID); err != nil {
			return fmt.Errorf("delete message: %w", err)
		}

		return nil
	}

	currentGame, err := vs.gameProvider.GetGameByID(ctx, game.IDFromString(currentPlayer.CurrentGameID))
	if err != nil {
		return fmt.Errorf("get game by id: %w", err)
	}

	playerAgainst := playerAgainstCurrentID(currentGame, currentPlayer.ID)

	shotStats, err := vs.playerStats.PlayerStats(ctx, playerAgainst.ID)
	if err != nil {
		return fmt.Errorf("player stats: %w", err)
	}

	if err := vs.notifier.StartGameWithStats(
		ctx,
		currentPlayer.ChatID,
		messageID,
		playerAgainst,
		shotStats,
	); err != nil {
		return fmt.Errorf("start game with stats notification: %w", err)
	}

	return nil
}

func playerAgainstCurrentID(currentGame game.Game, currentPlayerID player.ID) game.Player {
	if currentGame.State.Player1.ID == currentPlayerID {
		return currentGame.State.Player2
	}

	return currentGame.State.Player1
}
