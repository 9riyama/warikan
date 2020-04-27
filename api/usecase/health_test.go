package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/warikan/api/usecase"
)

func TestHealthUseCase_Execute(t *testing.T) {
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
			name:    "Repository err",
			mockErr: errors.New("repository error"),
			wantErr: usecase.InternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPingReposity{}
			mock.On("Ping").Return(tt.mockErr)

			u := usecase.NewHealthUseCase(mock)
			err := u.Execute()
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

type mockPingReposity struct {
	mock.Mock
}

func (m *mockPingReposity) Ping() error {
	return m.Called().Error(0)
}
