package infra

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/domain/repository"
	"github.com/warikan/api/infra/persistence"
)

func NewPaymentsRepository(db *sql.DB) *paymentPersistencePostgres {
	return &paymentPersistencePostgres{
		db: db,
	}
}

var _ repository.PaymentRepository = &paymentPersistencePostgres{}

type paymentPersistencePostgres struct {
	db *sql.DB
}

func (r *paymentPersistencePostgres) GetData(userID, cursor int) ([]*model.Payment, error) {
	const (
		limit = 20
	)

	payments, err := persistence.SelectPayments(r.db, userID, limit, cursor)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return payments, nil
}

func (r *paymentPersistencePostgres) Create(mp *model.Payment) (*model.Payment, error) {
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

func (*paymentPersistencePostgres) toModel(u *persistence.Payment) *model.Payment {
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

func (r *paymentPersistencePostgres) Update(mp *model.Payment) (*model.Payment, error) {
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

func (r *paymentPersistencePostgres) DeleteByID(userID, paymentID int) error {
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

func (r *paymentPersistencePostgres) FetchDate(userID int) ([]*string, error) {

	paymentsDate, err := persistence.SelectPaymentDateByUserID(r.db, userID)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return paymentsDate, nil
}
