package infra_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/infra"
	"github.com/warikan/db"
	"github.com/warikan/test"
)

func TestPaymentsPersistencePostgres_Create(t *testing.T) {

	r := infra.NewPaymentsRepository(db.Pool)

	now := time.Now()

	tests := []struct {
		name    string
		setup   func(t *testing.T)
		arg     *model.Payment
		want    *model.Payment
		wantErr error
	}{
		{
			name: "Success",
			setup: func(t *testing.T) {
				if err := test.LoadFixturesAt(db.Pool, "_fixtures"); err != nil {
					t.Fatal(err)
				}
			},
			arg: &model.Payment{
				UserID:      10001,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "作成", Valid: true},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			want: &model.Payment{
				UserID:      10001,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "作成", Valid: true},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got, err := r.Create(tt.arg)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, but got nil")
					return
				}

				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantErr) {
					t.Errorf("unexpected error type:\nwant: %T\ngot : %T", err, tt.wantErr)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("err should be nil, but got %q", err)
				return
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.Payment{}, "ID", "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPaymentsPersistencePostgres_Update(t *testing.T) {
	r := infra.NewPaymentsRepository(db.Pool)
	now := time.Now()

	tests := []struct {
		name    string
		setup   func(t *testing.T)
		arg     *model.Payment
		want    *model.Payment
		wantErr error
	}{
		{
			name: "Success",
			setup: func(t *testing.T) {
				if err := test.LoadFixturesAt(db.Pool, "_fixtures"); err != nil {
					t.Fatal(err)
				}
			},
			arg: &model.Payment{
				ID:          19999,
				UserID:      10001,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "更新後", Valid: true},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     5555,
			},
			want: &model.Payment{
				ID:          19999,
				UserID:      10001,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "更新後", Valid: true},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     5555,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Update(tt.arg)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("expected error, but got nil")
					return
				}

				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantErr) {
					t.Errorf("unexpected error type:\nwant: %T\ngot : %T", err, tt.wantErr)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("err should be nil, but got %q", err)
				return
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// func Test_PaymentRepository_DeleteByID(t *testing.T) {
// 	r := repository.NewPaymentsRepository(testDB)
// 	loadDefaultFixture(testDB, t)

// 	tests := []struct {
// 		name      string
// 		userID    int
// 		paymentID int
// 		wantErr   error
// 	}{
// 		{
// 			name:      "Success",
// 			userID:    1,
// 			paymentID: 1,
// 			wantErr:   nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := r.DeleteByID(tt.userID, tt.paymentID)
// 			if tt.wantErr != nil {
// 				if err == nil {
// 					t.Error("expected error, but got nil")
// 					return
// 				}

// 				err = errors.Cause(err)
// 				if g, e := err.Error(), tt.wantErr.Error(); g != e {
// 					t.Errorf("unexpected error:\nwant: %q\ngot : %q", e, g)
// 					return
// 				}
// 				return
// 			}
// 		})
// 	}
// }

// func loadDefaultFixture(db *sqlx.DB, t *testing.T) {
// 	tx := db.MustBegin()
// 	tx.MustExec(tx.Rebind("INSERT INTO users (user_name, partner_name, email, password, user_image, partner_image) VALUES (?, ?, ?, ?, ?, ?)"), "ユーザー", "パートナー", "example@co.jp", "password", "user_image", "partner_image")
// 	tx.MustExec(tx.Rebind("INSERT INTO categories (name) VALUES (?)"), "家賃")
// 	tx.MustExec(tx.Rebind("INSERT INTO categories (name) VALUES (?)"), "ガス代")
// 	tx.MustExec(tx.Rebind("INSERT INTO payers (name) VALUES (?)"), "ユーザー")
// 	tx.MustExec(tx.Rebind("INSERT INTO payers (name) VALUES (?)"), "パートナー")
// 	tx.MustExec(tx.Rebind("INSERT INTO payments (user_id, category_id, payer_id, payment_date, payment) VALUES (?, ?, ?, ?, ?)"), 1, 1, 1, "2020-04-01T00:00:00+09:00", 1234)
// 	_ = tx.Commit()
// }
