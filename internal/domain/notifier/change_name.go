package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) ChangeName(ctx context.Context, plr player.Player, fullName, userName string) error {
	//nolint:mnd
	buttons := make([]button.ButtonRow, 0, 2)

	buttons, err := appendChanageNameButton(buttons, fullName)
	if err != nil {
		return fmt.Errorf("create button by full name: %w", err)
	}

	buttons, err = appendChanageNameButton(buttons, userName)
	if err != nil {
		return fmt.Errorf("create button by username: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:  plr.ChatID,
		Text:    "Send message to change your name or use buttons for one of the your account name",
		Type:    msginfo.MessageTypePlain,
		Buttons: buttons,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func appendChanageNameButton(
	buttons []button.ButtonRow,
	displayName string,
) ([]button.ButtonRow, error) {
	if displayName != "" {
		btn, err := button.ChangeName(fmt.Sprintf("Use: %q", displayName), displayName)
		if err != nil {
			return nil, fmt.Errorf("create chanage name button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	return buttons, nil
}
