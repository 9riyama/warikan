package usecase_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/usecase"
)

func TestPaymentsUseCase_Create(t *testing.T) {
	tests := []struct {
		name    string
		param   *usecase.CreatePaymentParam
		userID  int
		mock    *model.Payment
		mockErr error
		want    *model.Payment
		wantErr error
	}{
		{
			name: "Success",
			param: &usecase.CreatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			userID: 1,
			mock: &model.Payment{
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			mockErr: nil,
			want: &model.Payment{
				ID:          0,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			wantErr: nil,
		},
		{
			name: "InvalidParam error",
			param: &usecase.CreatePaymentParam{
				Description: sql.NullString{String: "", Valid: false},
			},
			userID:  1,
			want:    nil,
			wantErr: usecase.InvalidParamError{},
		},
		{
			name: "Repository error",
			param: &usecase.CreatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			userID: 1,
			mock: &model.Payment{
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			mockErr: errors.New("repository error"),
			want:    nil,
			wantErr: usecase.InternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &mockPaymentRepository{}
			m.On("Create", tt.mock).Return(tt.want, tt.mockErr)

			u := usecase.NewPaymentUseCase(m)
			got, err := u.Create(tt.param, tt.userID)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, but got nil")
					return
				}
				if g, e := err.Error(), tt.wantErr.Error(); g != e {
					t.Errorf("unexpected error:\nwant: %v\ngot : %v", e, g)
				}
				return
			}

			if err != nil {
				t.Errorf("err should be nil, but got %q", err)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPaymentsUseCase_Update(t *testing.T) {
	tests := []struct {
		name      string
		param     *usecase.UpdatePaymentParam
		userID    int
		paymentID int
		mock      *model.Payment
		mockErr   error
		want      *model.Payment
		wantErr   error
	}{
		{
			name: "Success",
			param: &usecase.UpdatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			userID:    1,
			paymentID: 1,
			mock: &model.Payment{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			mockErr: nil,
			want: &model.Payment{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			wantErr: nil,
		},
		{
			name: "InvalidParam error",
			param: &usecase.UpdatePaymentParam{
				Description: sql.NullString{String: "", Valid: false},
			},
			userID:    1,
			paymentID: 1,
			want:      nil,
			wantErr:   usecase.InvalidParamError{},
		},
		{
			name: "Repository error",
			param: &usecase.UpdatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			userID:    1,
			paymentID: 1,
			mock: &model.Payment{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			mockErr: errors.New("repository error"),
			want:    nil,
			wantErr: usecase.InternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &mockPaymentRepository{}
			m.On("Update", tt.mock).Return(tt.want, tt.mockErr)

			u := usecase.NewPaymentUseCase(m)
			got, err := u.Update(tt.param, tt.userID, tt.paymentID)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, but got nil")
					return
				}
				if g, e := err.Error(), tt.wantErr.Error(); g != e {
					t.Errorf("unexpected error:\nwant: %v\ngot : %v", e, g)
				}
				return
			}

			if err != nil {
				t.Errorf("err should be nil, but got %q", err)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type mockPaymentRepository struct {
	mock.Mock
}

func (m *mockPaymentRepository) Create(mp *model.Payment) (*model.Payment, error) {
	ret := m.Called(mp)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}

func (m *mockPaymentRepository) Update(mp *model.Payment) (*model.Payment, error) {
	ret := m.Called(mp)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}
