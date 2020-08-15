package test

import (
	"database/sql"
	"path/filepath"

	testfixtures "github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
)

const (
	defaultPath = "/_fixtures"
)

func LoadFixtures(db *sql.DB) error {
	return LoadFixturesAt(db, defaultPath)
}

func LoadFixturesAt(db *sql.DB, fixturePath string) error {
	return LoadFixturesWithDBAt(db, fixturePath)
}

func LoadFixturesWithDBAt(db *sql.DB, path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return errors.WithStack(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(p),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = fixtures.Load(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
