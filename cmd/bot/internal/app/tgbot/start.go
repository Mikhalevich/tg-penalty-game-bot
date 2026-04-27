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

	http.ListenAndServe(":2000", t.bot.WebhookHandler())

	return nil
}
