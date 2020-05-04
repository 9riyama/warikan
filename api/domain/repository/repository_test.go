package repository_test

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	driver = "postgres"
	testDB *sqlx.DB
	err    error
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("WARIKAN_DB_DSN")
	testDB, err = sqlx.Connect(driver, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	code := m.Run()
	os.Exit(code)
}
