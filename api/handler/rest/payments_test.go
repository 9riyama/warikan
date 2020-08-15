package rest_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/go-cmp/cmp"
	mock "github.com/stretchr/testify/mock"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/handler/rest"
	"github.com/warikan/api/usecase"
)

func Test_paymentsHandler_GetData(t *testing.T) {
	tests := []struct {
		name         string
		strUserID    string
		strCursor    string
		payments     []*usecase.Payment
		userID       int
		cursor       int
		useCaseError error
		wantCode     int
		wantBody     string
	}{
		{
			name:      "Success",
			strUserID: "1",
			strCursor: "1",
			payments: []*usecase.Payment{
				{
					ID:           1,
					CategoryName: "カテゴリー名",
					PayerName:    "パートナー",
					PaymentDate:  "2020-04-01",
					Payment:      1234,
					CreatedAt:    "2020-04-01 09:00:00",
				},
			},
			userID:       1,
			cursor:       1,
			useCaseError: nil,
			wantCode:     200,
			wantBody:     "{\"payments\":[{\"id\":1,\"category_name\":\"カテゴリー名\",\"payer_name\":\"パートナー\",\"payment_date\":\"2020-04-01\",\"payment\":1234,\"created_at\":\"2020-04-01 09:00:00\"}]}\n",
		},
		{
			name:      "ConversionError",
			strUserID: "1",
			strCursor: "string",
			payments:  []*usecase.Payment{},
			userID:    999,
			cursor:    999,
			wantCode:  200,
			wantBody:  "{\"payments\":[]}\n",
		},
		{
			name:         "InternalServerError",
			strUserID:    "999",
			strCursor:    "999",
			payments:     []*usecase.Payment{},
			userID:       999,
			cursor:       999,
			useCaseError: usecase.InternalServerError{},
			wantCode:     500,
			wantBody:     `{"msg":"システム内部エラーが発生しました。"}` + "\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPaymentUseCase{}
			mock.On("GetData", tt.userID, tt.cursor).Return(tt.payments, tt.useCaseError)

			r := httptest.NewRequest(http.MethodGet, "/"+tt.strUserID+"/payments?cursor="+tt.strCursor, nil)
			rr := httptest.NewRecorder()
			h := rest.NewPaymentsHandler(mock)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tt.strUserID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.GetData(rr, r)

			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("GetData() mismatch status code (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantBody, rr.Body.String()); diff != "" {
				t.Errorf("GetData() mismatch body (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_paymentsHandler_CreateData(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		userID       string
		want         *model.Payment
		req          *usecase.CreatePaymentParam
		body         string
		useCaseError error
		wantCode     int
		wantBody     string
	}{
		{
			name:   "Success",
			id:     1,
			userID: "1",
			want: &model.Payment{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
				CreatedAt:   time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
			},
			req: &usecase.CreatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.Now().Location()),
				Payment:     1234,
			},
			body:         `{"category_id":1,"payer_id":1,"payment_date":"2020-04-01T00:00:00+09:00","payment":1234}`,
			useCaseError: nil,
			wantCode:     http.StatusCreated,
			wantBody:     "{\"id\":1,\"user_id\":1,\"category_id\":1,\"payer_id\":1,\"description\":{\"String\":\"\",\"Valid\":false},\"payment_date\":\"2020-04-01T00:00:00Z\",\"payment\":1234,\"created_at\":\"2020-04-01T00:00:00Z\",\"updated_at\":\"2020-04-01T00:00:00Z\"}\n",
		},
		{
			name:   "Internal server error",
			id:     1,
			userID: "1",
			want:   &model.Payment{},
			req: &usecase.CreatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.Now().Location()),
				Payment:     1234,
			},
			body:         `{"category_id":1,"payer_id":1,"payment_date":"2020-04-01T00:00:00+09:00","payment":1234}`,
			useCaseError: errors.New("usecase error"),
			wantCode:     http.StatusInternalServerError,
			wantBody:     `{"msg":"システム内部エラーが発生しました。"}` + "\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPaymentUseCase{}
			mock.On("Create", tt.req, tt.id).Return(tt.want, tt.useCaseError)

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			h := rest.NewPaymentsHandler(mock)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tt.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.CreateData(rr, r)

			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("CreateData() mismatch status code (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantBody, rr.Body.String()); diff != "" {
				t.Errorf("CreateData() mismatch body (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_paymentsHandler_UpdateData(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		pID          int
		paymentID    string
		userID       string
		want         *model.Payment
		req          *usecase.UpdatePaymentParam
		body         string
		useCaseError error
		wantCode     int
		wantBody     string
	}{
		{
			name:      "Success",
			id:        1,
			pID:       1,
			paymentID: "1",
			userID:    "1",
			want: &model.Payment{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
				CreatedAt:   time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
			},
			req: &usecase.UpdatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.Now().Location()),
				Payment:     1234,
			},
			body:         `{"category_id":1,"payer_id":1,"payment_date":"2020-04-01T00:00:00+09:00","payment":1234}`,
			useCaseError: nil,
			wantCode:     http.StatusOK,
			wantBody:     "{\"id\":1,\"user_id\":1,\"category_id\":1,\"payer_id\":1,\"description\":{\"String\":\"\",\"Valid\":false},\"payment_date\":\"2020-04-01T00:00:00Z\",\"payment\":1234,\"created_at\":\"2020-04-01T00:00:00Z\",\"updated_at\":\"2020-04-01T00:00:00Z\"}\n",
		},
		{
			name:      "Internal server error",
			id:        1,
			pID:       1,
			paymentID: "1",
			userID:    "1",
			want:      &model.Payment{},
			req: &usecase.UpdatePaymentParam{
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.Now().Location()),
				Payment:     1234,
			},
			body:         `{"category_id":1,"payer_id":1,"payment_date":"2020-04-01T00:00:00+09:00","payment":1234}`,
			useCaseError: errors.New("usecase error"),
			wantCode:     http.StatusInternalServerError,
			wantBody:     `{"msg":"システム内部エラーが発生しました。"}` + "\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPaymentUseCase{}
			mock.On("Update", tt.req, tt.id, tt.pID).Return(tt.want, tt.useCaseError)

			r := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			h := rest.NewPaymentsHandler(mock)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tt.userID)
			rctx.URLParams.Add("payment_id", tt.paymentID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.UpdateData(rr, r)

			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("CreateData() mismatch status code (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantBody, rr.Body.String()); diff != "" {
				t.Errorf("CreateData() mismatch body (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_paymentsHandler_DeleteData(t *testing.T) {
	tests := []struct {
		name         string
		strUserID    string
		userID       int
		strPaymentID string
		paymentID    int
		useCaseError error
		wantCode     int
		wantBody     string
	}{
		{
			name:         "Success",
			strUserID:    "1",
			userID:       1,
			strPaymentID: "1",
			paymentID:    1,
			useCaseError: nil,
			wantCode:     http.StatusNoContent,
			wantBody:     "",
		},
		{
			name:         "Bad request error paymentID is String",
			strUserID:    "1",
			userID:       1,
			strPaymentID: "string",
			paymentID:    1,
			useCaseError: errors.New("usecase error"),
			wantCode:     http.StatusBadRequest,
			wantBody:     `{"msg":"要求の形式が正しくありません。"}` + "\n",
		},
		{
			name:         "Bad request error userID is String",
			strUserID:    "string",
			userID:       1,
			strPaymentID: "1",
			paymentID:    1,
			useCaseError: errors.New("usecase error"),
			wantCode:     http.StatusBadRequest,
			wantBody:     `{"msg":"要求の形式が正しくありません。"}` + "\n",
		},
		{
			name:         "Internal server error",
			strUserID:    "1",
			userID:       1,
			strPaymentID: "1",
			paymentID:    1,
			useCaseError: errors.New("usecase error"),
			wantCode:     http.StatusInternalServerError,
			wantBody:     `{"msg":"システム内部エラーが発生しました。"}` + "\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &mockPaymentUseCase{}
			mock.On("DeleteByID", tt.userID, tt.paymentID).Return(tt.useCaseError)

			r := httptest.NewRequest(http.MethodDelete, "/", nil)
			rr := httptest.NewRecorder()
			h := rest.NewPaymentsHandler(mock)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tt.strUserID)
			rctx.URLParams.Add("payment_id", tt.strPaymentID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.DeleteData(rr, r)

			if diff := cmp.Diff(tt.wantCode, rr.Code); diff != "" {
				t.Errorf("DeleteData() mismatch status code (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantBody, rr.Body.String()); diff != "" {
				t.Errorf("DeleteData() mismatch body (-want +got):\n%s", diff)
			}
		})
	}
}

type mockPaymentUseCase struct {
	mock.Mock
	usecase.PaymentUseCase
}

func (m *mockPaymentUseCase) GetData(userID, cursor int) ([]*usecase.Payment, error) {
	ret := m.Called(userID, cursor)
	return ret.Get(0).([]*usecase.Payment), ret.Error(1)
}

func (m *mockPaymentUseCase) Create(param *usecase.CreatePaymentParam, userID int) (*model.Payment, error) {
	ret := m.Called(param, userID)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}

func (m *mockPaymentUseCase) Update(param *usecase.UpdatePaymentParam, userID, paymentID int) (*model.Payment, error) {
	ret := m.Called(param, userID, paymentID)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}

func (m *mockPaymentUseCase) DeleteByID(userID, paymentID int) error {
	ret := m.Called(userID, paymentID)
	return ret.Error(0)
}
