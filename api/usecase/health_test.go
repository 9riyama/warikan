package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/warikan/api/domain/repository"
	"github.com/warikan/api/usecase"
)

func Test_healthUseCase_Check(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr error
	}{
		{
			name:    "Success",
			mockErr: nil,
			wantErr: nil,
		},
		{
			name:    "Repository error",
			mockErr: errors.New("repository error"),
			wantErr: usecase.ServiceUnavailableError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPingRepository{}
			mock.On("Ping").Return(tt.mockErr)

			u := usecase.NewHealthUseCase(mock)
			err := u.Check(context.Background())
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("err should be nil, but got %q", err)
				}
				return
			}

			if err == nil {
				t.Error("expected error, but got nil")
				return
			}
			if g, e := err.Error(), tt.wantErr.Error(); g != e {
				t.Errorf("unexpected error:\nwant: %q\ngot : %q", e, g)
			}
		})
	}
}

var _ repository.PingRepository = &mockPingRepository{}

type mockPingRepository struct {
	mock.Mock
}

func (m *mockPingRepository) Ping(ctx context.Context) error {
	return m.Called().Error(0)
}
