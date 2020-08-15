package usecase_test

import (
	"os"
	"testing"

	"github.com/warikan/log"
)

func TestMain(m *testing.M) {
	log.Init()
	code := m.Run()
	os.Exit(code)
}
