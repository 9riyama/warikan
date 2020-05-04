package usecase

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/domain/repository"
)

type PaymentUseCase interface {
	Create(req *CreatePaymentParam, userID int) (*model.Payment, error)
}

func NewPaymentUseCase(r repository.PaymentRepository) *paymentUsecase {
	return &paymentUsecase{r}
}

var _ PaymentUseCase = &paymentUsecase{}

type paymentUsecase struct {
	PaymentRepository repository.PaymentRepository
}

type CreatePaymentParam struct {
	CategoryID  int            `json:"category_id" validate:"required"`
	PayerID     int            `json:"payer_id" validate:"required"`
	Description sql.NullString `json:"description"`
	PaymentDate time.Time      `json:"payment_date" validate:"required"`
	Payment     int            `json:"payment" validate:"required"`
}

func (u *paymentUsecase) Create(param *CreatePaymentParam, userID int) (*model.Payment, error) {

	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		log.Println("validation error")
		return nil, InvalidParamError{}
	}

	payment := &model.Payment{
		UserID:      userID,
		CategoryID:  param.CategoryID,
		PayerID:     param.PayerID,
		Description: param.Description,
		PaymentDate: param.PaymentDate,
		Payment:     param.Payment,
	}

	payment, err = u.PaymentRepository.Create(payment)

	if err != nil {
		log.Println("repository error")
		return nil, InternalServerError{}
	}
	return payment, nil
}