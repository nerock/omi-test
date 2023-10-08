package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nerock/omi-test/account/internal"
)

type AccountPostgresStorage struct {
	db *sql.DB
}

func NewAccountPostgresStorage(db *sql.DB) *AccountPostgresStorage {
	return &AccountPostgresStorage{db: db}
}

func (as *AccountPostgresStorage) UpdateAccountByID(ctx context.Context, a internal.Account) error {
	const queryUpdateAccountNameByID = `
UPDATE
	accounts
SET
	name=$2
WHERE
    id = $1
`

	res, err := as.db.ExecContext(ctx, queryUpdateAccountNameByID, a.ID, a.Name)
	if err != nil {
		return fmt.Errorf("could not execute update: %w", err)
	}

	if rows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("invalid result: %w", err)
	} else if rows != 1 {
		return internal.ErrInvalidAccount
	}

	return nil
}
