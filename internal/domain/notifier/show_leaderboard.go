package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
)

const (
	leaderboardMaxButtonsCount = 2
)

func (n *Notifier) ShowLeaderboard(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	playerID player.ID,
	pageNumber int,
	pagesCount int,
	positions []player.Position,
) error {
	buttons, err := makeLeaderboardButtons(pageNumber, pagesCount)
	if err != nil {
		return fmt.Errorf("make leaderboard buttons: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     chatID,
		ReplyMsgID: messageID,
		Text:       n.makeLeaderboardMsg(playerID, positions),
		Type:       sendOrEditTextMarkdownType(messageID),
		Buttons:    buttons,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeLeaderboardButtons(pageNumber, pagesCount int) ([]button.ButtonRow, error) {
	row := make(button.ButtonRow, 0, leaderboardMaxButtonsCount)

	if pageNumber > 1 {
		left, err := button.LeaderboardPage("<", pageNumber-1)
		if err != nil {
			return nil, fmt.Errorf("create left button: %w", err)
		}

		row = append(row, left)
	}

	if pageNumber < pagesCount {
		right, err := button.LeaderboardPage(">", pageNumber+1)
		if err != nil {
			return nil, fmt.Errorf("create right button: %w", err)
		}

		row = append(row, right)
	}

	return []button.ButtonRow{row}, nil
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
