package notifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
)

func (n *Notifier) StartGameWithStats(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	playerAgainst game.Player,
	shotStats game.ShotStats,
) error {
	if err := n.sendStartGameWithStats(
		ctx,
		chatID,
		messageID,
		n.statsMsg(playerAgainst, shotStats),
	); err != nil {
		return fmt.Errorf("send start game to first player: %w", err)
	}

	return nil
}

func (n *Notifier) statsMsg(plr game.Player, shotStats game.ShotStats) string {
	return fmt.Sprintf("%s\n*Attack*\n%s\n*Defend*\n%s",
		n.startGameAgainstMsg(plr),
		n.playerStatsMsg(shotStats.Attack),
		n.playerStatsMsg(shotStats.Defend),
	)
}

func (n *Notifier) playerStatsMsg(shotStats []game.ShotSidePercent) string {
	if len(shotStats) == 0 {
		return "No shots"
	}

	lines := make([]string, 0, len(shotStats))

	for _, ss := range shotStats {
		lines = append(lines,
			n.escaper.EscapeMarkdown(
				fmt.Sprintf("%s %5.2f%%",
					symbolBySide(ss.Side),
					ss.Percent)),
		)
	}

	return strings.Join(lines, "\n")
}

func symbolBySide(side game.ShotSide) string {
	switch side {
	case game.ShotSideLeft:
		return leftSideSymbol

	case game.ShotSideRight:
		return rightSideSymbol

	case game.ShotSideMiddle:
		return middleSideSymbol

	case game.ShotSideMiss:
		return missSymbol

	case game.ShotSideNoShot:
	}

	return side.String()
}

func (n *Notifier) sendStartGameWithStats(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	msg string,
) error {
	if err := n.sender.SendMessage(ctx, msginfo.Message{
		ChatID:     chatID,
		ReplyMsgID: messageID,
		Text:       msg,
		Type:       msginfo.MessageTypeEditMarkdown,
		Buttons: []button.ButtonRow{
			{
				game.LeaveGameButton("Leave"),
			},
		},
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
