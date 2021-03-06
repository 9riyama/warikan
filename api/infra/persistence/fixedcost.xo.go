// Package persistence contains the types for schema 'public'.
package persistence

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"time"
)

// FixedCost represents a row from 'public.fixed_costs'.
type FixedCost struct {
	ID          int            `json:"id"`           // id
	UserID      int            `json:"user_id"`      // user_id
	CategoryID  int            `json:"category_id"`  // category_id
	PayerID     int            `json:"payer_id"`     // payer_id
	Description sql.NullString `json:"description"`  // description
	PaymentDate time.Time      `json:"payment_date"` // payment_date
	Payment     int            `json:"payment"`      // payment
	CreatedAt   time.Time      `json:"created_at"`   // created_at
	UpdatedAt   time.Time      `json:"updated_at"`   // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the FixedCost exists in the database.
func (fc *FixedCost) Exists() bool {
	return fc._exists
}

// Deleted provides information if the FixedCost has been deleted from the database.
func (fc *FixedCost) Deleted() bool {
	return fc._deleted
}

// Insert inserts the FixedCost to the database.
func (fc *FixedCost) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if fc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.fixed_costs (` +
		`user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt)
	err = db.QueryRow(sqlstr, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt).Scan(&fc.ID)
	if err != nil {
		return err
	}

	// set existence
	fc._exists = true

	return nil
}

// Update updates the FixedCost in the database.
func (fc *FixedCost) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if fc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.fixed_costs SET (` +
		`user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) WHERE id = $9`

	// run query
	XOLog(sqlstr, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt, fc.ID)
	_, err = db.Exec(sqlstr, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt, fc.ID)
	return err
}

// Save saves the FixedCost to the database.
func (fc *FixedCost) Save(db XODB) error {
	if fc.Exists() {
		return fc.Update(db)
	}

	return fc.Insert(db)
}

// Upsert performs an upsert for FixedCost.
//
// NOTE: PostgreSQL 9.5+ only
func (fc *FixedCost) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if fc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.fixed_costs (` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.user_id, EXCLUDED.category_id, EXCLUDED.payer_id, EXCLUDED.description, EXCLUDED.payment_date, EXCLUDED.payment, EXCLUDED.created_at, EXCLUDED.updated_at` +
		`)`

	// run query
	XOLog(sqlstr, fc.ID, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt)
	_, err = db.Exec(sqlstr, fc.ID, fc.UserID, fc.CategoryID, fc.PayerID, fc.Description, fc.PaymentDate, fc.Payment, fc.CreatedAt, fc.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	fc._exists = true

	return nil
}

// Delete deletes the FixedCost from the database.
func (fc *FixedCost) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !fc._exists {
		return nil
	}

	// if deleted, bail
	if fc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.fixed_costs WHERE id = $1`

	// run query
	XOLog(sqlstr, fc.ID)
	_, err = db.Exec(sqlstr, fc.ID)
	if err != nil {
		return err
	}

	// set deleted
	fc._deleted = true

	return nil
}

// Category returns the Category associated with the FixedCost's CategoryID (category_id).
//
// Generated from foreign key 'fixed_costs_category_id_fkey'.
func (fc *FixedCost) Category(db XODB) (*Category, error) {
	return CategoryByID(db, fc.CategoryID)
}

// Payer returns the Payer associated with the FixedCost's PayerID (payer_id).
//
// Generated from foreign key 'fixed_costs_payer_id_fkey'.
func (fc *FixedCost) Payer(db XODB) (*Payer, error) {
	return PayerByID(db, fc.PayerID)
}

// User returns the User associated with the FixedCost's UserID (user_id).
//
// Generated from foreign key 'fixed_costs_user_id_fkey'.
func (fc *FixedCost) User(db XODB) (*User, error) {
	return UserByID(db, fc.UserID)
}

// FixedCostsByCategoryID retrieves a row from 'public.fixed_costs' as a FixedCost.
//
// Generated from index 'fixed_costs_category_id_idx'.
func FixedCostsByCategoryID(db XODB, categoryID int) ([]*FixedCost, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at ` +
		`FROM public.fixed_costs ` +
		`WHERE category_id = $1`

	// run query
	XOLog(sqlstr, categoryID)
	q, err := db.Query(sqlstr, categoryID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*FixedCost{}
	for q.Next() {
		fc := FixedCost{
			_exists: true,
		}

		// scan
		err = q.Scan(&fc.ID, &fc.UserID, &fc.CategoryID, &fc.PayerID, &fc.Description, &fc.PaymentDate, &fc.Payment, &fc.CreatedAt, &fc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &fc)
	}

	return res, nil
}

// FixedCostsByPayerID retrieves a row from 'public.fixed_costs' as a FixedCost.
//
// Generated from index 'fixed_costs_payer_id_idx'.
func FixedCostsByPayerID(db XODB, payerID int) ([]*FixedCost, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at ` +
		`FROM public.fixed_costs ` +
		`WHERE payer_id = $1`

	// run query
	XOLog(sqlstr, payerID)
	q, err := db.Query(sqlstr, payerID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*FixedCost{}
	for q.Next() {
		fc := FixedCost{
			_exists: true,
		}

		// scan
		err = q.Scan(&fc.ID, &fc.UserID, &fc.CategoryID, &fc.PayerID, &fc.Description, &fc.PaymentDate, &fc.Payment, &fc.CreatedAt, &fc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &fc)
	}

	return res, nil
}

// FixedCostByID retrieves a row from 'public.fixed_costs' as a FixedCost.
//
// Generated from index 'fixed_costs_pkey'.
func FixedCostByID(db XODB, id int) (*FixedCost, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at ` +
		`FROM public.fixed_costs ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	fc := FixedCost{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&fc.ID, &fc.UserID, &fc.CategoryID, &fc.PayerID, &fc.Description, &fc.PaymentDate, &fc.Payment, &fc.CreatedAt, &fc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}

// FixedCostsByUserID retrieves a row from 'public.fixed_costs' as a FixedCost.
//
// Generated from index 'fixed_costs_user_id_idx'.
func FixedCostsByUserID(db XODB, userID int) ([]*FixedCost, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, category_id, payer_id, description, payment_date, payment, created_at, updated_at ` +
		`FROM public.fixed_costs ` +
		`WHERE user_id = $1`

	// run query
	XOLog(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*FixedCost{}
	for q.Next() {
		fc := FixedCost{
			_exists: true,
		}

		// scan
		err = q.Scan(&fc.ID, &fc.UserID, &fc.CategoryID, &fc.PayerID, &fc.Description, &fc.PaymentDate, &fc.Payment, &fc.CreatedAt, &fc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &fc)
	}

	return res, nil
}
