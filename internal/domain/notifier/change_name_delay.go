package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ChangeNameDelay(ctx context.Context, plr player.Player, waitPeriod time.Duration) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: plr.ChatID,
		Text:   fmt.Sprintf("You can change name only after *%s*", waitPeriod),
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send too many requests message: %w", err)
	}

	return nil
}
