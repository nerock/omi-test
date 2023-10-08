package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type mockAuditStorage struct {
	saveLogFunc func(ctx context.Context, log AuditLog) error
}

func (m *mockAuditStorage) SaveLog(ctx context.Context, log AuditLog) error {
	return m.saveLogFunc(ctx, log)
}

func TestAuditService_AddLog(t *testing.T) {
	Convey("AddLog", t, func() {
		ctx := context.Background()
		st := &mockAuditStorage{}
		svc := NewAuditService(st)

		newLog := AuditLog{
			EventTypeID: 1,
			EventType:   "ACCOUNT_UPDATED",
			Timestamp:   time.Now(),
			UserIP:      "192.168.1.1",
			Data:        []byte(`"account_id": 123`),
		}
		Convey("when storage returns an error", func() {
			st.saveLogFunc = func(inCtx context.Context, inLog AuditLog) error {
				So(inCtx, ShouldResemble, ctx)
				So(inLog, ShouldResemble, newLog)

				return errors.New("something happened")
			}

			Convey("the service should return an error", func() {
				err := svc.AddLog(ctx, newLog)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("when storage returns no error", func() {
			st.saveLogFunc = func(inCtx context.Context, inLog AuditLog) error {
				So(inCtx, ShouldResemble, ctx)
				So(inLog, ShouldResemble, newLog)

				return nil
			}

			Convey("the service should return no error", func() {
				err := svc.AddLog(ctx, newLog)
				So(err, ShouldBeNil)
			})
		})
	})
}
