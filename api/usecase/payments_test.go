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

func Test_paymentUsecase_GetData(t *testing.T) {

	tests := []struct {
		name     string
		userID   int
		cursor   int
		mockWant []*model.Payment
		mockErr  error
		want     []*usecase.Payment
		wantErr  error
	}{
		{
			name:   "Success",
			userID: 1,
			cursor: 1,
			mockWant: []*model.Payment{
				{
					ID:           1,
					CategoryName: "カテゴリー名",
					PayerName:    "パートナー",
					PaymentDate:  time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
					Payment:      1234,
					CreatedAt:    time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockErr: nil,
			want: []*usecase.Payment{
				{
					ID:           1,
					CategoryName: "カテゴリー名",
					PayerName:    "パートナー",
					PaymentDate:  "2020-04-01",
					Payment:      1234,
					CreatedAt:    "2020-04-01 09:00:00",
				},
			},
			wantErr: nil,
		},
		{
			name:     "Repository error",
			userID:   1,
			cursor:   1,
			mockWant: []*model.Payment{},
			mockErr:  errors.New("repository error"),
			want:     []*usecase.Payment{},
			wantErr:  usecase.InternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPaymentRepository{}
			mock.On("GetData", tt.userID, tt.cursor).Return(tt.mockWant, tt.mockErr)

			u := usecase.NewPaymentUseCase(mock)
			got, err := u.GetData(tt.userID, tt.cursor)
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
				t.Error("expected error, but got nil")
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("usecase returned wrong response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

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

func TestPaymentsUseCase_DeleteByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		paymentID int
		mockErr   error
		wantErr   error
	}{
		{
			name:      "Success",
			userID:    1,
			paymentID: 1,
			mockErr:   nil,
			wantErr:   nil,
		},
		{
			name:      "Repository error",
			userID:    1,
			paymentID: 1,
			mockErr:   errors.New("repository error"),
			wantErr:   usecase.InternalServerError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &mockPaymentRepository{}
			m.On("DeleteByID", tt.userID, tt.paymentID).Return(tt.mockErr)

			u := usecase.NewPaymentUseCase(m)
			err := u.DeleteByID(tt.userID, tt.paymentID)
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
		})
	}
}

type mockPaymentRepository struct {
	mock.Mock
}

func (m *mockPaymentRepository) GetData(userID, cursor int) ([]*model.Payment, error) {
	ret := m.Called(userID, cursor)
	return ret.Get(0).([]*model.Payment), ret.Error(1)
}

func (m *mockPaymentRepository) Create(mp *model.Payment) (*model.Payment, error) {
	ret := m.Called(mp)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}

func (m *mockPaymentRepository) Update(mp *model.Payment) (*model.Payment, error) {
	ret := m.Called(mp)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}

func (m *mockPaymentRepository) DeleteByID(userID, paymentID int) error {
	ret := m.Called(userID, paymentID)
	return ret.Error(0)
}

func (m *mockPaymentRepository) FetchDate(userID int) ([]*model.Payment, error) {
	ret := m.Called(userID)
	return ret.Get(0).([]*model.Payment), ret.Error(1)
}
