package internal

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidAccount = errors.New("invalid account")

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
	if id == "" || name == "" {
		return Account{}, ErrInvalidAccount
	}

	acc := Account{
		ID:   id,
		Name: name,
	}
	if err := a.st.UpdateAccountByID(ctx, acc); err != nil {
		return Account{}, fmt.Errorf("could not update account: %w", err)
	}

	return acc, nil
}
