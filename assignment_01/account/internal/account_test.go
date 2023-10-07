package internal

import (
	"context"
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

	})
}
