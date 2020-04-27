package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	mock "github.com/stretchr/testify/mock"
	"github.com/warikan/api/handler/rest"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		want    int
	}{
		{
			name:    "Success",
			mockErr: nil,
			want:    http.StatusOK,
		},
		{
			name:    "Use case error",
			mockErr: errors.New("use case error"),
			want:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := new(mockHealthUseCase)
			mock.On("Execute").Return(tt.mockErr)

			r := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			h := rest.NewHealthHandler(mock)
			h.Health(rr, r)
			if diff := cmp.Diff(tt.want, rr.Code); diff != "" {
				t.Errorf("Health() mismatch status code (-want +got):\n%s", diff)
			}
		})
	}
}

type mockHealthUseCase struct {
	mock.Mock
}

func (m *mockHealthUseCase) Execute() error {
	ret := m.Called()
	return ret.Error(0)
}
