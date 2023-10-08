package storage

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nerock/omi-test/logger/internal"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuditPostgresStorage_SaveLog(t *testing.T) {
	Convey("SaveLog", t, func() {
		ctx := context.Background()
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		st := NewAuditPostgresStorage(db)

		query := "INSERT INTO audit_logs (type, timestamp, user_ip, data) VALUES ($1, $2, $3, $4)"
		log := internal.AuditLog{
			EventTypeID: 10,
			EventType:   "ACTION_HAPPENED",
			Timestamp:   time.Now(),
			UserIP:      "192.168.1.1",
			Data:        []byte(`"changes": "theChange"`),
		}
		Convey("when the database returns an error", func() {
			mock.ExpectExec(query).
				WillReturnError(errors.New("something happened"))

			Convey("the storage should return an error", func() {
				So(st.SaveLog(ctx, log), ShouldNotBeNil)
			})
		})

		Convey("when the database returns an invalid result", func() {
			mock.ExpectExec(query).
				WillReturnResult(sqlmock.NewResult(0, 0))

			Convey("the storage should return an invalid account error", func() {
				So(st.SaveLog(ctx, log), ShouldNotBeNil)
			})
		})

		Convey("when the database returns a valid result", func() {
			mock.ExpectExec(query).
				WillReturnResult(sqlmock.NewResult(0, 1))

			Convey("the storage should return no error", func() {
				So(st.SaveLog(ctx, log), ShouldBeNil)
			})
		})
	})
}
