package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	nameButtonsCount = 3
)

func (n *Notifier) ChangeName(
	ctx context.Context,
	plr player.Player,
	fullName,
	userName string,
) error {
	buttons, err := makeChanageNameButton(fullName, userName)
	if err != nil {
		return fmt.Errorf("create button by full name: %w", err)
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

func makeChanageNameButton(
	fullName,
	userName string,
) ([]button.ButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, nameButtonsCount)

	if fullName != "" {
		btn, err := button.ChangeName(fmt.Sprintf("Use: %q", fullName), fullName)
		if err != nil {
			return nil, fmt.Errorf("create chanage name button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	if userName != "" {
		btn, err := button.ChangeName(fmt.Sprintf("Use: %q", userName), userName)
		if err != nil {
			return nil, fmt.Errorf("create chanage name button: %w", err)
		}

		buttons = append(buttons, button.Row(btn))
	}

	buttons = append(buttons, button.Row(button.ChangeNameCancel("Cancel")))

	return buttons, nil
}
