package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/warikan/api/usecase"
)

type PaymentsHandler interface {
	GetData(http.ResponseWriter, *http.Request)
	CreateData(http.ResponseWriter, *http.Request)
	UpdateData(http.ResponseWriter, *http.Request)
	DeleteData(http.ResponseWriter, *http.Request)
	FetchDate(http.ResponseWriter, *http.Request)
}

type paymentsHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentsHandler(u usecase.PaymentUseCase) PaymentsHandler {
	return &paymentsHandler{
		useCase: u,
	}
}

type paymentHandlerResponse struct {
	Payments []*usecase.Payment `json:"payments"`
}

type paymentsDateResponse struct {
	PaymentsDate []*usecase.PaymentDate `json:"payments_date"`
}

func (h *paymentsHandler) GetData(w http.ResponseWriter, r *http.Request) {

	strCursor := r.URL.Query().Get("cursor")
	if strCursor == "" {
		strCursor = "0"
	}
	cursor, err := strconv.Atoi(strCursor)

	strUserID := chi.URLParam(r, "user_id")
	userID, _ := strconv.Atoi(strUserID)

	res := paymentHandlerResponse{}
	if err != nil {
		res.Payments = make([]*usecase.Payment, 0)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			internalServerError(w)
		}
		return
	}

	payments, err := h.useCase.GetData(userID, cursor)
	if err != nil {
		httpError(w, err)
		return
	}

	res.Payments = payments

	if err := json.NewEncoder(w).Encode(res); err != nil {
		internalServerError(w)
	}
}

func (h *paymentsHandler) CreateData(w http.ResponseWriter, r *http.Request) {
	strUserID := chi.URLParam(r, "user_id")
	userID, _ := strconv.Atoi(strUserID)

	req := usecase.CreatePaymentParam{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequestError(w)
		return
	}

	resp, err := h.useCase.Create(&req, userID)
	if err != nil {
		httpError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpError(w, err)
	}
}

func (h *paymentsHandler) UpdateData(w http.ResponseWriter, r *http.Request) {
	strUserID := chi.URLParam(r, "user_id")
	strPayemntID := chi.URLParam(r, "payment_id")
	userID, _ := strconv.Atoi(strUserID)
	payemntID, _ := strconv.Atoi(strPayemntID)

	req := usecase.UpdatePaymentParam{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequestError(w)
		return
	}

	resp, err := h.useCase.Update(&req, userID, payemntID)
	if err != nil {
		httpError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpError(w, err)
	}
}

func (h *paymentsHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	strUserID := chi.URLParam(r, "user_id")
	strPayemntID := chi.URLParam(r, "payment_id")

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		badRequestError(w)
		return
	}
	payemntID, err := strconv.Atoi(strPayemntID)
	if err != nil {
		badRequestError(w)
		return
	}

	if err := h.useCase.DeleteByID(userID, payemntID); err != nil {
		httpError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *paymentsHandler) FetchDate(w http.ResponseWriter, r *http.Request) {

	strUserID := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		badRequestError(w)
		return
	}

	paymentsDate, err := h.useCase.FetchDate(userID)
	if err != nil {
		httpError(w, err)
		return
	}
	res := paymentsDateResponse{PaymentsDate: paymentsDate}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		internalServerError(w)
	}
}
