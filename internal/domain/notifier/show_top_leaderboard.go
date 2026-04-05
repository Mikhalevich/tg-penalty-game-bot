package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ShowTopLeaderbord(
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
				button.LeaderboardPlayer("My position"),
			},
		},
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
