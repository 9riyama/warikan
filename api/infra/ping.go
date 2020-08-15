package infra

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/warikan/api/domain/repository"
)

func NewPingPersistencePostgres(db *sql.DB) *pingPersistencePostgres {
	return &pingPersistencePostgres{
		db: db,
	}
}

var _ repository.PingRepository = &pingPersistencePostgres{}

type pingPersistencePostgres struct {
	db *sql.DB
}

func (p *pingPersistencePostgres) Ping(ctx context.Context) error {
	return errors.WithStack(p.db.PingContext(ctx))
}
