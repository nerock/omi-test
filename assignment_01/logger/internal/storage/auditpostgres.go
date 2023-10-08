package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/nerock/omi-test/logger/internal"
)

type jsonb []byte

func (j *jsonb) Value() (driver.Value, error) {
	if len(*j) == 0 {
		return nil, nil
	}

	return string(*j), nil
}

func (j *jsonb) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}

	s, ok := value.([]byte)
	if !ok {
		errors.New("scan source was not a byte slice")
	}

	copy(*j, s)

	return nil
}

type AuditPostgresStorage struct {
	db *sql.DB
}

func NewAuditPostgresStorage(db *sql.DB) *AuditPostgresStorage {
	return &AuditPostgresStorage{db: db}
}

func (as *AuditPostgresStorage) SaveLog(ctx context.Context, log internal.AuditLog) error {
	const query = `
INSERT INTO
	audit_logs (type, timestamp, user_ip, data)
VALUES
    ($1, $2, $3, $4)`

	res, err := as.db.ExecContext(ctx, query, log.EventType, log.Timestamp.Unix(), log.UserIP, jsonb(log.Data))
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	if rows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("could not check rows: %w", err)
	} else if rows == 0 {
		return fmt.Errorf("no logs added")
	}

	return nil
}
