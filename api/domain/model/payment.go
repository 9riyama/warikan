package model

import (
	"database/sql"
	"time"
)

type Payment struct {
	ID          int            `json:"id"`
	UserID      int            `json:"user_id"`
	CategoryID  int            `json:"category_id"`
	PayerID     int            `json:"payer_id"`
	Description sql.NullString `json:"description"`
	PaymentDate time.Time      `json:"payment_date"`
	Payment     int            `json:"payment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
