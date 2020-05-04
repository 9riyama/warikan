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
			wantBody:     "{\"message\":\"Internal Server Error\"}\n",
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
			rctx.URLParams.Add("id", tt.userID)
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

type mockPaymentUseCase struct {
	mock.Mock
	usecase.PaymentUseCase
}

func (m *mockPaymentUseCase) Create(req *usecase.CreatePaymentParam, userID int) (*model.Payment, error) {
	ret := m.Called(req, userID)
	return ret.Get(0).(*model.Payment), ret.Error(1)
}
