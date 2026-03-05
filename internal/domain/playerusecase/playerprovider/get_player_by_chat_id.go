package playerprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/player"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/infra/logger"
)

func (p *PlayerProvider) GetPlayerByChatID(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	plr, err := p.repo.GetPlayerByChatID(ctx, chatID)
	if err != nil {
		if !p.repo.IsNotFoundError(err) {
			return player.Player{}, fmt.Errorf("get player by chat_id: %w", err)
		}

		newPlayer, err := p.createPlayer(ctx, chatID)
		if err != nil {
			return player.Player{}, fmt.Errorf("create new player: %w", err)
		}

		return newPlayer, nil
	}

	return plr, nil
}

func (p *PlayerProvider) createPlayer(
	ctx context.Context,
	chatID msginfo.ChatID,
) (player.Player, error) {
	var (
		creationTime = p.timeProvider.Now()
	)

	for {
		newPlayer := player.Player{
			ChatID:      chatID,
			DisplayName: p.nameGenerator.GenerateName(),
			CreatedAt:   creationTime,
			GameStatus:  player.GameStatusIdle,
		}

		playerID, err := p.repo.CreatePlayer(ctx, newPlayer)
		if err != nil {
			if p.repo.IsAlreadyExistsError(err) {
				logger.FromContext(ctx).
					WithField("player_name", newPlayer.DisplayName).
					Warn("trying to grab existing name")

				continue
			}

			return player.Player{}, fmt.Errorf("create new player: %w", err)
		}

		newPlayer.ID = playerID

		return newPlayer, nil
	}
}
