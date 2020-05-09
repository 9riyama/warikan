package persistence

import (
	"github.com/warikan/api/domain/model"
)

func SelectPayments(db XODB, userID, limit, cursor int) ([]*model.Payment, error) {
	var err error

	args := []interface{}{userID, limit}

	// sql query
	var sqlstr = `SELECT p.id
		, c.name AS category_name
		, a.name AS payer_name
		, p.payment_date
		, p.payment
		, p.created_at
		FROM payments p
		LEFT JOIN payers a
		ON p.payer_id = a.id
		LEFT JOIN categories c
		ON p.category_id = c.id
		WHERE p.user_id = $1`

	if cursor != 0 {
		sqlstr += `WHERE a.id < $3`
		args = append(args, cursor)
	}

	sqlstr += `
		ORDER BY p.payment_date, p.created_at DESC
		LIMIT $2`

	// run query
	XOLog(sqlstr, args...)
	q, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}

	defer q.Close()

	payments := make([]*model.Payment, 0)
	for q.Next() {
		var p model.Payment
		err := q.Scan(
			&p.ID,
			&p.CategoryName,
			&p.PayerName,
			&p.PaymentDate,
			&p.Payment,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}

	return payments, nil
}
