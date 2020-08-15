package rest

import (
	"encoding/json"
	"net/http"

	"github.com/warikan/api/usecase"
)

func NewHealthHandler(uc usecase.HealthUseCase, version string) *HealthHandler {
	return &HealthHandler{
		uc:      uc,
		version: version,
	}
}

type HealthHandler struct {
	uc      usecase.HealthUseCase
	version string
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := struct {
		Status  healthStatus `json:"status"`
		Version string       `json:"version"`
	}{
		Status:  pass,
		Version: h.version,
	}

	if err := h.uc.Check(ctx); err != nil {
		resp.Status = fail
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpError(w, err, "")
	}
}

type healthStatus string

const (
	pass healthStatus = "pass"
	fail healthStatus = "fail"
	//warn healthStatus = "warn"
)
