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
	msg, isContaintsCurrentPlayer := n.makeLeaderboardMsg(playerID, positions)

	buttons, err := makeLeaderboardButtons(pageNumber, pagesCount, isContaintsCurrentPlayer)
	if err != nil {
		return fmt.Errorf("make leaderboard buttons: %w", err)
	}

	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     chatID,
		ReplyMsgID: messageID,
		Text:       msg,
		Type:       sendOrEditTextMarkdownType(messageID),
		Buttons:    buttons,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeLeaderboardButtons(
	pageNumber int,
	pagesCount int,
	isContainsCurrentPlayer bool,
) ([]button.ButtonRow, error) {
	navigationRow, err := makeLeaderboardNavigationButtonRow(pageNumber, pagesCount)
	if err != nil {
		return nil, fmt.Errorf("make navigation buttons: %w", err)
	}

	if !isContainsCurrentPlayer {
		return []button.ButtonRow{
			navigationRow,
			button.Row(button.LeaderboardMyPage("My position")),
		}, nil
	}

	if pageNumber == 1 {
		return []button.ButtonRow{
			navigationRow,
		}, nil
	}

	topBtn, err := button.LeaderboardPage("Top", 1)
	if err != nil {
		return nil, fmt.Errorf("create top button: %w", err)
	}

	return []button.ButtonRow{
		navigationRow,
		button.Row(topBtn),
	}, nil
}

func makeLeaderboardNavigationButtonRow(
	pageNumber int,
	pagesCount int,
) (button.ButtonRow, error) {
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

	if len(row) == 0 {
		return nil, nil
	}

	return row, nil
}

func sendOrEditTextMarkdownType(messageID msginfo.MessageID) msginfo.MessageType {
	if messageID.Int() == 0 {
		return msginfo.MessageTypeMarkdown
	}

	return msginfo.MessageTypeEditMarkdown
}

// makeLeaderboardMsg construct leaderboard message
// returns constructed message and flag that shows current player in positions.
func (n *Notifier) makeLeaderboardMsg(
	playerID player.ID,
	positions []player.Position,
) (string, bool) {
	var (
		lines                   = make([]string, 0, len(positions))
		isContainsCurrentPlayer = false
	)

	for _, pos := range positions {
		if playerID == pos.ID {
			lines = append(lines, fmt.Sprintf("%d\\. *%s* %d",
				pos.Position,
				n.escaper.EscapeMarkdown(pos.DisplayName),
				pos.Score,
			))

			isContainsCurrentPlayer = true

			continue
		}

		lines = append(lines, fmt.Sprintf("%d\\. %s %d",
			pos.Position,
			n.escaper.EscapeMarkdown(pos.DisplayName),
			pos.Score,
		))
	}

	return strings.Join(lines, "\n"), isContainsCurrentPlayer
}
