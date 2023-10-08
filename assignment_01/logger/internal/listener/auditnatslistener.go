package listener

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nerock/omi-test/logger/internal"
	"github.com/nerock/omi-test/logger/pb/event"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"time"
)

type AuditNatsListener struct {
	c   *nats.Conn
	svc AuditService
	log *zap.Logger
}

type AuditService interface {
	AddLog(ctx context.Context, log internal.AuditLog) error
}

func Subscribe(ctx context.Context, c *nats.Conn, topic string, svc AuditService, log *zap.Logger) {
	c.Subscribe(topic, func(msg *nats.Msg) {
		data := &event.AuditEvent{}
		if err := proto.Unmarshal(msg.Data, data); err != nil {
			log.Error("could not unmarshal event data", zap.Error(err))
			msg.Nak()
			return
		}

		timeStamp, err := time.Parse(time.RFC3339, data.GetTimestamp())
		if err != nil {
			log.Error("could not parse timestamp", zap.Error(err))
			msg.Nak()
			return
		}

		eventData, err := json.Marshal(data.GetData())
		if err != nil {
			log.Error("could not serialize event data", zap.Error(err))
			msg.Nak()
			return
		}

		if err := svc.AddLog(ctx, internal.AuditLog{
			EventTypeID: int(data.GetEventType()),
			EventType:   data.EventType.String(),
			Timestamp:   timeStamp,
			UserIP:      data.GetUserIp(),
			Data:        eventData,
		}); err != nil {
			log.Error("could not add log", zap.Error(err))
			msg.Nak()
			return
		}

		msg.Ack()
	})
}
