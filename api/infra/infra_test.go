package infra_test

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"github.com/warikan/db"
	"github.com/warikan/log"
)

var (
	maxconn  = 5
	filePath = "../../_config/test.config.yaml"
)

func TestMain(m *testing.M) {
	log.Init()

	if err := db.Init(maxconn, filePath); err != nil {
		log.Logger.Fatal("failed to initDB:", zap.Error(err))
	}
	defer db.Close()

	code := m.Run()
	os.Exit(code)
}
