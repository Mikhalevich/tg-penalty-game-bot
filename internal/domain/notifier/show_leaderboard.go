package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ShowLeaderbord(
	ctx context.Context,
	chatID msginfo.ChatID,
	playerID player.ID,
	positions []player.Position,
) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   n.makeLeaderboardMsg(playerID, positions),
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (n *Notifier) makeLeaderboardMsg(playerID player.ID, positions []player.Position) string {
	lines := make([]string, 0, len(positions))

	for _, pos := range positions {
		if playerID == pos.ID {
			lines = append(lines, fmt.Sprintf("%d\\. *%s* %d",
				pos.Position,
				n.escaper.EscapeMarkdown(pos.DisplayName),
				pos.Score,
			))

			continue
		}

		lines = append(lines, fmt.Sprintf("%d\\. %s %d",
			pos.Position,
			n.escaper.EscapeMarkdown(pos.DisplayName),
			pos.Score,
		))
	}

	return strings.Join(lines, "\n")
}
