package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ShowPlayerLeaderbord(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	playerID player.ID,
	positions []player.Position,
) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     chatID,
		ReplyMsgID: messageID,
		Text:       n.makeLeaderboardMsg(playerID, positions),
		Type:       sendOrEditTextMarkdownType(messageID),
		Buttons: []button.ButtonRow{
			{
				button.LeaderboardTop("Top"),
			},
		},
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func sendOrEditTextMarkdownType(messageID msginfo.MessageID) msginfo.MessageType {
	if messageID.Int() == 0 {
		return msginfo.MessageTypeMarkdown
	}

	return msginfo.MessageTypeEditMarkdown
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
