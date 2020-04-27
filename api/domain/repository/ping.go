package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type HealthRepository interface {
	Ping() error
}

func NewHealthRepository(db *sqlx.DB) *healthRepository {
	return &healthRepository{db}
}

var _ HealthRepository = &healthRepository{}

type healthRepository struct {
	db *sqlx.DB
}

func (r *healthRepository) Ping() error {
	const sql = "SELECT 1"
	_, err := r.db.Exec(sql)
	return errors.WithStack(err)
}
