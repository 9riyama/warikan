package log_test

import (
	"os"
	"testing"

	"github.com/warikan/log"
)

func TestInit(t *testing.T) {
	const envKey = "GO_ENV"

	tests := []struct {
		name string
		env  string
	}{
		{
			name: "Production logger",
			env:  "production",
		},
		{
			name: "Default logger",
			env:  "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() func() {
				e := os.Getenv(envKey)
				_ = os.Setenv(envKey, tt.env)
				return func() {
					_ = os.Setenv(envKey, e)
				}
			}()()

			log.Init()
			if log.Logger == nil {
				t.Error("Logger is nil")
			}
		})
	}
}
