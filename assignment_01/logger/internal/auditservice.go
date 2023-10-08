package internal

import (
	"context"
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

	return nil
}
