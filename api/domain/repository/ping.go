package repository

import (
	"context"
)

type PingRepository interface {
	Ping(ctx context.Context) error
}
