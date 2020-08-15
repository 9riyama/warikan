package db

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/warikan/config"
	"github.com/warikan/log"
)

var (
	Pool *sql.DB
)

func Init(maxconn int, filePath string) error {
	dsn, err := config.GetDSN(filePath)
	if err != nil {
		return errors.WithStack(err)
	}

	Pool, err = sql.Open("pgx", dsn)
	if err != nil {
		return errors.WithStack(err)
	}
	Pool.SetMaxIdleConns(maxconn)
	Pool.SetMaxOpenConns(maxconn)
	Pool.SetConnMaxLifetime(10 * time.Second)

	log.Logger.Debug("database connections initialized", zap.Int("maxconn", maxconn))

	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
