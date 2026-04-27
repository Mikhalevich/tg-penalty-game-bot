package tgbot

import (
	"context"
	"fmt"
	"net/http"
)

func (t *TGBot) Start(ctx context.Context) error {
	if err := t.setMyCommands(ctx); err != nil {
		return fmt.Errorf("set my commands: %w", err)
	}

	if !t.isWebHook {
		t.bot.Start(ctx)

		return nil
	}

	go t.bot.StartWebhook(ctx)

	//nolint:gosec
	if err := http.ListenAndServe(":2000", t.bot.WebhookHandler()); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
