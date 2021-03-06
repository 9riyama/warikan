// Package persistence contains the types for schema 'public'.
package persistence

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"

	"github.com/lib/pq"
)

// User represents a row from 'public.users'.
type User struct {
	ID           int         `json:"id"`            // id
	UserName     string      `json:"user_name"`     // user_name
	PartnerName  string      `json:"partner_name"`  // partner_name
	Email        string      `json:"email"`         // email
	Password     string      `json:"password"`      // password
	UserImage    string      `json:"user_image"`    // user_image
	PartnerImage string      `json:"partner_image"` // partner_image
	Proportion   int16       `json:"proportion"`    // proportion
	CreatedAt    time.Time   `json:"created_at"`    // created_at
	UpdatedAt    time.Time   `json:"updated_at"`    // updated_at
	DeletedAt    pq.NullTime `json:"deleted_at"`    // deleted_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted provides information if the User has been deleted from the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.users (` +
		`user_name, partner_name, email, password, user_image, partner_image, proportion, created_at, updated_at, deleted_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
	err = db.QueryRow(sqlstr, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt).Scan(&u.ID)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// Update updates the User in the database.
func (u *User) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if u._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.users SET (` +
		`user_name, partner_name, email, password, user_image, partner_image, proportion, created_at, updated_at, deleted_at` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10` +
		`) WHERE id = $11`

	// run query
	XOLog(sqlstr, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt, u.ID)
	_, err = db.Exec(sqlstr, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt, u.ID)
	return err
}

// Save saves the User to the database.
func (u *User) Save(db XODB) error {
	if u.Exists() {
		return u.Update(db)
	}

	return u.Insert(db)
}

// Upsert performs an upsert for User.
//
// NOTE: PostgreSQL 9.5+ only
func (u *User) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.users (` +
		`id, user_name, partner_name, email, password, user_image, partner_image, proportion, created_at, updated_at, deleted_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, user_name, partner_name, email, password, user_image, partner_image, proportion, created_at, updated_at, deleted_at` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.user_name, EXCLUDED.partner_name, EXCLUDED.email, EXCLUDED.password, EXCLUDED.user_image, EXCLUDED.partner_image, EXCLUDED.proportion, EXCLUDED.created_at, EXCLUDED.updated_at, EXCLUDED.deleted_at` +
		`)`

	// run query
	XOLog(sqlstr, u.ID, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
	_, err = db.Exec(sqlstr, u.ID, u.UserName, u.PartnerName, u.Email, u.Password, u.UserImage, u.PartnerImage, u.Proportion, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// Delete deletes the User from the database.
func (u *User) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.users WHERE id = $1`

	// run query
	XOLog(sqlstr, u.ID)
	_, err = db.Exec(sqlstr, u.ID)
	if err != nil {
		return err
	}

	// set deleted
	u._deleted = true

	return nil
}

// UserByID retrieves a row from 'public.users' as a User.
//
// Generated from index 'users_pkey'.
func UserByID(db XODB, id int) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_name, partner_name, email, password, user_image, partner_image, proportion, created_at, updated_at, deleted_at ` +
		`FROM public.users ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&u.ID, &u.UserName, &u.PartnerName, &u.Email, &u.Password, &u.UserImage, &u.PartnerImage, &u.Proportion, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
