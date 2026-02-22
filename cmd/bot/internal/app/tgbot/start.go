package tgbot

import (
	"context"
	"fmt"
)

func (t *TGBot) Start(ctx context.Context) error {
	if err := t.setMyCommands(ctx); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	t.bot.Start(ctx)

	return nil
}
