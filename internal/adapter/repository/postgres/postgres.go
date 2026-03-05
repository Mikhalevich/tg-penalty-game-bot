package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gamecontroller"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller/changename"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller/changestatus"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller/findgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller/playerprovider"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playercontroller/startgame"
)

var (
	_ changename.Repository     = (*Postgres)(nil)
	_ changestatus.Repository   = (*Postgres)(nil)
	_ findgame.Repository       = (*Postgres)(nil)
	_ playerprovider.Repository = (*Postgres)(nil)
	_ startgame.Repository      = (*Postgres)(nil)
	_ gamecontroller.Repository = (*Postgres)(nil)
)

type Driver interface {
	IsConstraintError(err error, constraint string) bool
}

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
	ExtContext(ctx context.Context) sqlx.ExtContext
}

type Postgres struct {
	dbDriver   Driver
	transactor Transactor
}

func New(
	dbDriver Driver,
	transactor Transactor,
) *Postgres {
	return &Postgres{
		transactor: transactor,
		dbDriver:   dbDriver,
	}
}

func (p *Postgres) Transactor() Transactor {
	return p.transactor
}
