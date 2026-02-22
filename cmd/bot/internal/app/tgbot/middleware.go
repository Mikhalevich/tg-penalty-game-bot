package tgbot

type Middleware func(next Handler) Handler

func (t *TGBot) MiddlewareGroup(next func(tbot *TGBot)) {
	group := &TGBot{
		bot:         t.bot,
		logger:      t.logger,
		middlewares: t.middlewares[:len(t.middlewares):len(t.middlewares)],
	}

	next(group)
}

func (t *TGBot) AddMiddleware(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

func (t *TGBot) applyMiddleware(h Handler) Handler {
	for i := len(t.middlewares) - 1; i >= 0; i-- {
		h = t.middlewares[i](h)
	}

	return h
}
