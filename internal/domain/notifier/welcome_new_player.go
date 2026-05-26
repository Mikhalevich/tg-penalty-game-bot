package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

func (n *Notifier) WelcomeNewPlayer(ctx context.Context, plr player.Player) error {
	welcomeMsg := fmt.Sprintf("Welcome to the penalty shootout game bot\nYour name is *%s*",
		n.escaper.EscapeMarkdown(plr.DisplayName))

	playBotBtn, err := game.StartGameButton("Play with bot", game.GameTypeBot)
	if err != nil {
		return fmt.Errorf("create bot button: %w", err)
	}

	playByLinkBtn, err := game.StartGameButton("Play by share link", game.GameTypeByLink)
	if err != nil {
		return fmt.Errorf("create play by link button: %w", err)
	}

	playRatingBtn, err := game.StartGameButton("Play rating game", game.GameTypeRating)
	if err != nil {
		return fmt.Errorf("create rating button: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID: plr.ChatID,
		Text:   welcomeMsg,
		Type:   msginfo.MessageTypeMarkdown,
		Buttons: []button.ButtonRow{
			button.Row(playBotBtn),
			button.Row(playByLinkBtn),
			button.Row(playRatingBtn),
			button.Row(button.ChangeNameTrigger("Change name")),
		},
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
