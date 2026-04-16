package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/forceshot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/gameshot"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/gameshotstate"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/gameusecase/startgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/matchmaking"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/notifier"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/outboxprocessor"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/changename"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/changestatus"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/findgame"
	"github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/playerprovider"
	playerstartgame "github.com/Mikhalevich/tg-penalty-game-bot/internal/domain/playerusecase/startgame"
)

var (
	_ changename.Repository      = (*Postgres)(nil)
	_ changestatus.Repository    = (*Postgres)(nil)
	_ findgame.Repository        = (*Postgres)(nil)
	_ playerprovider.Repository  = (*Postgres)(nil)
	_ playerstartgame.Repository = (*Postgres)(nil)

	_ matchmaking.Repository = (*Postgres)(nil)

	_ gameshot.Repository      = (*Postgres)(nil)
	_ forceshot.Repository     = (*Postgres)(nil)
	_ startgame.Repository     = (*Postgres)(nil)
	_ gameshotstate.Repository = (*Postgres)(nil)

	_ outboxprocessor.Repository = (*Postgres)(nil)

	_ notifier.Sender = (*Postgres)(nil)
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
