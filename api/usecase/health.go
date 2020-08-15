package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/warikan/api/domain/repository"
	"github.com/warikan/log"
)

type HealthUseCase interface {
	Check(ctx context.Context) error
}

func NewHealthUseCase(pr repository.PingRepository) *healthUseCase {
	return &healthUseCase{
		pr: pr,
	}
}

var _ HealthUseCase = &healthUseCase{}

type healthUseCase struct {
	pr repository.PingRepository
}

func (uc *healthUseCase) Check(ctx context.Context) error {
	if err := uc.pr.Ping(ctx); err != nil {
		log.Logger.Error("unhealthy database:", zap.Error(err))
		return ServiceUnavailableError{}
	}
	return nil
}
