package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) ShowLeaderboardRestrict(ctx context.Context, chatID msginfo.ChatID) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: chatID,
		Text:   "You need to play rating game to show leaderboard",
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
