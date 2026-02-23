package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) NameChanged(ctx context.Context, plr player.Player, msgID msginfo.MessageID) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     plr.ChatID,
		ReplyMsgID: msgID,
		Text:       fmt.Sprintf("Name changed to *%s*", n.escaper.EscapeMarkdown(plr.DisplayName)),
		Type:       msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
