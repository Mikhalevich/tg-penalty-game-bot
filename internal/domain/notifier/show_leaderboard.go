package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ShowLeaderbord(ctx context.Context, chatID msginfo.ChatID, positions []player.Position) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   n.makeLeaderboardMsg(positions),
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) makeLeaderboardMsg(positions []player.Position) string {
	lines := make([]string, 0, len(positions))

	for _, pos := range positions {
		lines = append(lines, fmt.Sprintf("%d\\. %s %d",
			pos.Position,
			n.escaper.EscapeMarkdown(pos.DisplayName),
			pos.Score,
		))
	}

	return strings.Join(lines, "\n")
}
