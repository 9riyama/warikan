package repository

import (
	"github.com/warikan/api/domain/model"
)

type PaymentRepository interface {
	GetData(userID, cursor int) ([]*model.Payment, error)
	Create(*model.Payment) (*model.Payment, error)
	Update(*model.Payment) (*model.Payment, error)
	DeleteByID(userID, paymentID int) error
	FetchDate(userID int) ([]*string, error)
}
