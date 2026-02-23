package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) WelcomeNewPlayer(ctx context.Context, plr player.Player) error {
	welcomeMsg := fmt.Sprintf("Welcome to the penalty shootout game bot\nYour name is *%s*",
		n.escaper.EscapeMarkdown(plr.DisplayName))

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: plr.ChatID,
		Text:   welcomeMsg,
		Type:   msginfo.MessageTypeMarkdown,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
