package startgame

import (
	"context"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/game"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/perror"
)

func (s *StartGame) RepeatGame(
	ctx context.Context,
	chatID msginfo.ChatID,
	gameType game.GameType,
) error {
	switch gameType {
	case game.GameTypeRating:
		return s.FindGame(ctx, chatID)

	case game.GameTypeBot:
		return s.PlayWithBot(ctx, chatID)

	case game.GameTypeByLink:
		return s.StartGameByLink(ctx, chatID)
	}

	return perror.InvalidParam("invalid game type")
}
