package rest

import (
	"net/http"

	"github.com/warikan/api/usecase"
)

type HealthHandler interface {
	Health(http.ResponseWriter, *http.Request)
}

func NewHealthHandler(u usecase.HealthUseCase) *healthHandler {
	return &healthHandler{
		useCase: u,
	}
}

type healthHandler struct {
	useCase usecase.HealthUseCase
}

func (hu healthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if err := hu.useCase.Execute(); err != nil {
		httpError(w, err)
	}
}
