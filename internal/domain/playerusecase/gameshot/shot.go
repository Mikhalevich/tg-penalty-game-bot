package gameshot

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (g *GameShot) Shot(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	gameID game.ID,
	round int,
	side game.ShotSide,
) error {
	currentPlayer, err := g.playerProvider.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get player: %w", err)
	}

	if !currentPlayer.IsInGame(gameID.String()) {
		if err := g.notifier.PlayerNotInGame(ctx, chatID); err != nil {
			return fmt.Errorf("not in game notification: %w", err)
		}

		return nil
	}

	if err := g.gameController.Shot(ctx, game.Shot{
		GameID:       gameID,
		PlayerID:     currentPlayer.ID,
		Round:        round,
		ExpectedSide: side,
		CompletedAt:  g.timeProvider.Now(),
	}); err != nil {
		return fmt.Errorf("game shot: %w", err)
	}

	return nil
}
