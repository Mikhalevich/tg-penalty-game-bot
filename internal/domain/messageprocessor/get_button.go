package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/model/button"
)

func (m *MessageProcessor) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	btn, err := m.buttonRepository.GetButton(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get button from repository: %w", err)
	}

	return btn, nil
}
