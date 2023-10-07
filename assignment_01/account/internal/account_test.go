package internal

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockStorage struct {
	updateAccountByIDFunc func(context.Context, Account) error
}

func (m *mockStorage) UpdateAccountByID(ctx context.Context, a Account) error {
	return m.updateAccountByIDFunc(ctx, a)
}

func TestAccountService_UpdateAccount(t *testing.T) {
	Convey("UpdateAccount", t, func() {
		ctx := context.Background()
		st := &mockStorage{}
		svc := NewAccountService(st)

		Convey("when receiving invalid input", func() {
			var id, name string

			Convey("the service should return an empty account and a validation error", func() {
				a, err := svc.UpdateAccount(ctx, id, name)

				So(a, ShouldBeZeroValue)
				So(err, ShouldBeError, ErrInvalidAccount)
			})
		})

		Convey("when receiving valid input", func() {
			id, name := "id", "name"

			Convey("and the storage fails", func() {
				st.updateAccountByIDFunc = func(inCtx context.Context, inAccount Account) error {
					So(inCtx, ShouldResemble, ctx)
					So(inAccount, ShouldResemble, Account{ID: id, Name: name})

					return errors.New("something happened")
				}

				Convey("the service should return an empty account and an error", func() {
					a, err := svc.UpdateAccount(ctx, id, name)

					So(a, ShouldBeZeroValue)
					So(err, ShouldNotBeNil)
				})
			})

			Convey("and the storage succeeds", func() {
				st.updateAccountByIDFunc = func(inCtx context.Context, inAccount Account) error {
					So(inCtx, ShouldResemble, ctx)
					So(inAccount, ShouldResemble, Account{ID: id, Name: name})

					return nil
				}

				Convey("the service should return the updated account and no error", func() {
					a, err := svc.UpdateAccount(ctx, id, name)

					So(a, ShouldResemble, Account{ID: id, Name: name})
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
