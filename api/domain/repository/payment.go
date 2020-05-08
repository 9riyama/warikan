package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/infra/persistence"
)

type PaymentRepository interface {
	Create(*model.Payment) (*model.Payment, error)
	Update(*model.Payment) (*model.Payment, error)
	DeleteByID(userID, paymentID int) error
}

func NewPaymentRepository(db *sqlx.DB) *paymentRepository {
	return &paymentRepository{db}
}

var _ PaymentRepository = &paymentRepository{}

type paymentRepository struct {
	db *sqlx.DB
}

func (r *paymentRepository) Create(mp *model.Payment) (*model.Payment, error) {
	now := time.Now()

	p := &persistence.Payment{
		UserID:      mp.UserID,
		CategoryID:  mp.CategoryID,
		PayerID:     mp.PayerID,
		Description: mp.Description,
		PaymentDate: mp.PaymentDate,
		Payment:     mp.Payment,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := p.Save(r.db); err != nil {
		return nil, errors.WithStack(err)
	}

	payment := r.toModel(p)

	return payment, nil
}

func (*paymentRepository) toModel(u *persistence.Payment) *model.Payment {
	payment := &model.Payment{
		ID:          u.ID,
		UserID:      u.UserID,
		CategoryID:  u.CategoryID,
		PayerID:     u.PayerID,
		Description: u.Description,
		PaymentDate: u.PaymentDate,
		Payment:     u.Payment,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	return payment
}

func (r *paymentRepository) Update(mp *model.Payment) (*model.Payment, error) {
	now := time.Now()

	p, err := persistence.PaymentByID(r.db, mp.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	p.CategoryID = mp.CategoryID
	p.PayerID = mp.PayerID
	p.Description = mp.Description
	p.PaymentDate = mp.PaymentDate
	p.Payment = mp.Payment
	p.UpdatedAt = now

	if err := p.Save(r.db); err != nil {
		return nil, errors.WithStack(err)
	}

	payment := r.toModel(p)

	return payment, nil
}

func (r *paymentRepository) DeleteByID(userID, paymentID int) error {
	p, err := persistence.PaymentByID(r.db, paymentID)
	if err != nil {
		return errors.WithStack(err)
	}

	if p.UserID != userID {
		return errors.WithStack(err)
	}

	err = p.Delete(r.db)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
