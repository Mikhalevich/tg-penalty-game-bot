package postgres

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNoRowsUpdated = errors.New("no rows updated")
)

func (p *Postgres) IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func (p *Postgres) IsAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

func (p *Postgres) IsNoRowsUpdated(err error) bool {
	return errors.Is(err, ErrNoRowsUpdated)
}
