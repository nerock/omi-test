package storage

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nerock/omi-test/account/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAccountPostgresStorage_UpdateAccountByID(t *testing.T) {
	Convey("UpdateAccountByID", t, func() {
		ctx := context.Background()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		st := NewAccountPostgresStorage(db)

		query := "UPDATE accounts SET name=$2 WHERE id = $1"
		acc := internal.Account{ID: "id", Name: "name"}
		Convey("when the database returns an error", func() {
			mock.ExpectExec(query).
				WillReturnError(errors.New("something happened"))

			Convey("the storage should return an error", func() {
				So(st.UpdateAccountByID(ctx, acc), ShouldNotBeNil)
			})
		})

		Convey("when the database returns an invalid result", func() {
			mock.ExpectExec(query).
				WillReturnResult(sqlmock.NewResult(0, 0))

			Convey("the storage should return an invalid account error", func() {
				So(st.UpdateAccountByID(ctx, acc), ShouldBeError, internal.ErrInvalidAccount)
			})
		})

		Convey("when the database returns a valid result", func() {
			mock.ExpectExec(query).
				WillReturnResult(sqlmock.NewResult(0, 1))

			Convey("the storage should return no error", func() {
				So(st.UpdateAccountByID(ctx, acc), ShouldBeNil)
			})
		})
	})
}
