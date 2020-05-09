package usecase

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/domain/repository"
	"github.com/warikan/api/usecase/util"
)

type PaymentUseCase interface {
	GetData(userID, cursor int) ([]*Payment, error)
	Create(req *CreatePaymentParam, userID int) (*model.Payment, error)
	Update(req *UpdatePaymentParam, userID int, paymentID int) (*model.Payment, error)
	DeleteByID(userID, paymentID int) error
}

func NewPaymentUseCase(r repository.PaymentRepository) *paymentUsecase {
	return &paymentUsecase{r}
}

var _ PaymentUseCase = &paymentUsecase{}

type paymentUsecase struct {
	PaymentRepository repository.PaymentRepository
}

type Payment struct {
	ID           int    `json:"id"`
	CategoryName string `json:"category_name"`
	PayerName    string `json:"payer_name"`
	PaymentDate  string `json:"payment_date"`
	Payment      int    `json:"payment"`
	CreatedAt    string `json:"created_at"`
}

type CreatePaymentParam struct {
	CategoryID  int            `json:"category_id" validate:"required"`
	PayerID     int            `json:"payer_id" validate:"required"`
	Description sql.NullString `json:"description"`
	PaymentDate time.Time      `json:"payment_date" validate:"required"`
	Payment     int            `json:"payment" validate:"required"`
}

type UpdatePaymentParam struct {
	ID          int            `json:"-"`
	CategoryID  int            `json:"category_id" validate:"required"`
	PayerID     int            `json:"payer_id" validate:"required"`
	Description sql.NullString `json:"description"`
	PaymentDate time.Time      `json:"payment_date" validate:"required"`
	Payment     int            `json:"payment" validate:"required"`
}

func (u *paymentUsecase) GetData(userID, cursor int) ([]*Payment, error) {

	p, err := u.PaymentRepository.GetData(userID, cursor)
	if err != nil {
		log.Println("internal server error")
		return nil, InternalServerError{}
	}

	payments := make([]*Payment, 0, len(p))

	for _, v := range p {

		res := &Payment{
			ID:           v.ID,
			CategoryName: v.CategoryName,
			PayerName:    v.PayerName,
			PaymentDate:  util.ConvertJSTStringDate(v.PaymentDate),
			Payment:      v.Payment,
			CreatedAt:    util.ConvertJSTStringTime(v.CreatedAt),
		}
		payments = append(payments, res)

	}

	return payments, nil
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

func (u *paymentUsecase) Update(param *UpdatePaymentParam, userID, paymentID int) (*model.Payment, error) {

	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		log.Println("validation error")
		return nil, InvalidParamError{}
	}

	payment := &model.Payment{
		ID:          paymentID,
		UserID:      userID,
		CategoryID:  param.CategoryID,
		PayerID:     param.PayerID,
		Description: param.Description,
		PaymentDate: param.PaymentDate,
		Payment:     param.Payment,
	}

	payment, err = u.PaymentRepository.Update(payment)

	if err != nil {
		log.Println("repository error")
		return nil, InternalServerError{}
	}
	return payment, nil
}

func (u *paymentUsecase) DeleteByID(userID, paymentID int) error {
	err := u.PaymentRepository.DeleteByID(userID, paymentID)
	if err != nil {
		log.Println("repository error")
		return InternalServerError{}
	}
	return nil
}
