package storage

import (
	"context"
	"database/sql"
	"github.com/nerock/omi-test/account/internal"
)

type AccountPostgresStorage struct {
	db *sql.DB
}

func NewAccountPostgresStorage(db *sql.DB) *AccountPostgresStorage {
	return &AccountPostgresStorage{db: db}
}

func (as *AccountPostgresStorage) UpdateAccountByID(ctx context.Context, a internal.Account) error {
	return nil
}
