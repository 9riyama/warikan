package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"

	"github.com/warikan/api/handler/rest"
	"github.com/warikan/api/usecase"
)

func TestHealthHandler_Check(t *testing.T) {
	tests := []struct {
		name     string
		mockErr  error
		wantCode int
		wantBody string
	}{
		{
			name:     "Success",
			mockErr:  nil,
			wantCode: http.StatusOK,
			wantBody: `{"status":"pass","version":"unknown"}` + "\n",
		},
		{
			name:     "Unhealthy",
			mockErr:  usecase.ServiceUnavailableError{},
			wantCode: http.StatusServiceUnavailable,
			wantBody: `{"status":"fail","version":"unknown"}` + "\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockHealthUseCase{}
			mock.On("Check").Return(tt.mockErr)

			h := rest.NewHealthHandler(mock, "unknown")
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			h.Check(rr, r)
			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("Check() mismatch status code (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("Check() mismatch body (-want +got):\n%s", diff)
			}
		})
	}
}

var _ usecase.HealthUseCase = &mockHealthUseCase{}

type mockHealthUseCase struct {
	mock.Mock
}

func (m *mockHealthUseCase) Check(ctx context.Context) error {
	return m.Called().Error(0)
}
