package repository_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/warikan/api/domain/model"
	"github.com/warikan/api/domain/repository"
)

func Test_UserRepository_Create(t *testing.T) {
	r := repository.NewPaymentRepository(testDB)
	now := time.Now()

	tests := []struct {
		name    string
		arg     *model.Payment
		setup   func()
		want    *model.Payment
		wantErr error
	}{
		{
			name: "Success",
			arg: &model.Payment{
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
				PaymentDate: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC),
				Payment:     1234,
			},
			setup: func() {
				loadDefaultFixture(testDB, t)
			},
			want: &model.Payment{
				ID:          0,
				UserID:      1,
				CategoryID:  1,
				PayerID:     1,
				Description: sql.NullString{String: "", Valid: false},
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
			tt.setup()
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

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(model.Payment{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func loadDefaultFixture(db *sqlx.DB, t *testing.T) {
	tx := db.MustBegin()
	tx.MustExec(tx.Rebind("INSERT INTO users (user_name, partner_name, email, password, user_image, partner_image) VALUES (?, ?, ?, ?, ?, ?)"), "ユーザー", "パートナー", "example@co.jp", "password", "user_image", "partner_image")
	tx.MustExec(tx.Rebind("INSERT INTO categories (name) VALUES (?)"), "家賃")
	tx.MustExec(tx.Rebind("INSERT INTO payers (name) VALUES (?)"), "ユーザー")
	tx.MustExec(tx.Rebind("INSERT INTO payers (name) VALUES (?)"), "パートナー")
	_ = tx.Commit()
}
