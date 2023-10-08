package internal

import (
	"context"
	"fmt"
	"time"
)

type Storage interface {
	SaveLog(ctx context.Context, log AuditLog) error
}

type AuditLog struct {
	EventTypeID int
	EventType   string
	Timestamp   time.Time
	UserIP      string
	Data        []byte
}

type AuditService struct {
	st Storage
}

func NewAuditService(st Storage) *AuditService {
	return &AuditService{st: st}
}

func (a *AuditService) AddLog(ctx context.Context, log AuditLog) error {
	if err := a.st.SaveLog(ctx, log); err != nil {
		return fmt.Errorf("could not add new audit log: %w", err)
	}

	return nil
}
