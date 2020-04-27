package usecase

import (
	"log"

	"github.com/warikan/api/domain/repository"
)

type HealthUseCase interface {
	Execute() error
}

func NewHealthUseCase(r repository.HealthRepository) *healthUseCase {
	return &healthUseCase{
		repository: r,
	}
}

type healthUseCase struct {
	repository repository.HealthRepository
}

func (h *healthUseCase) Execute() error {
	err := h.repository.Ping()
	if err != nil {
		log.Printf("failed to ping to database%v\n", err)
		return InternalServerError{}
	}
	return nil
}
