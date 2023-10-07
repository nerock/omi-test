package internal

import "context"

type Account struct {
	ID   string
	Name string
}

type storage interface {
	UpdateAccountByID(ctx context.Context, a Account) error
}

type AccountService interface {
	UpdateAccount(ctx context.Context, id, name string) (Account, error)
}

type accountService struct {
	st storage
}

func NewAccountService(st storage) *accountService {
	return &accountService{st: st}
}

func (a *accountService) UpdateAccount(ctx context.Context, id, name string) (Account, error) {
	return Account{}, nil
}
