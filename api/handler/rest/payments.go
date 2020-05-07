package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/warikan/api/usecase"
)

type PaymentsHandler interface {
	CreateData(http.ResponseWriter, *http.Request)
	UpdateData(http.ResponseWriter, *http.Request)
	DeleteData(http.ResponseWriter, *http.Request)
}

type paymentsHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentsHandler(u usecase.PaymentUseCase) PaymentsHandler {
	return &paymentsHandler{
		useCase: u,
	}
}

func (h *paymentsHandler) CreateData(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(i)

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
	strUserID := chi.URLParam(r, "id")
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

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		httpError(w, err)
	}
}

func (h *paymentsHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	strPayemntID := chi.URLParam(r, "payment_id")
	payemntID, err := strconv.Atoi(strPayemntID)
	if err != nil {
		badRequestError(w)
		return
	}

	if err := h.useCase.DeleteByID(payemntID); err != nil {
		httpError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
