package interceptor

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nerock/omi-test/account/internal"
	"github.com/nerock/omi-test/account/pb/event"
	"go.uber.org/zap"
	"google.golang.org/grpc/peer"
	"time"
)

type AccountNatsInterceptor struct {
	c     *nats.EncodedConn
	topic string
	svc   internal.AccountService
	log   *zap.Logger
}

func NewAccountNatsInterceptor(svc internal.AccountService, c *nats.EncodedConn, topic string, log *zap.Logger) *AccountNatsInterceptor {
	return &AccountNatsInterceptor{
		svc:   svc,
		c:     c,
		topic: topic,
		log:   log,
	}
}

func (a *AccountNatsInterceptor) UpdateAccount(ctx context.Context, id, name string) (internal.Account, error) {
	var ip string
	peer, ok := peer.FromContext(ctx)
	if ok {
		ip = peer.Addr.String()
	}

	acc, err := a.svc.UpdateAccount(ctx, id, name)
	if err == nil {
		if err := a.c.Publish(a.topic, &event.AuditEvent{
			EventType: event.EventType_ACCOUNT_UPDATED,
			Timestamp: time.Now().Format(time.RFC3339),
			UserIp:    ip,
			Data:      &event.AuditEvent_Account{Account: &event.AccountData{AccountId: id}},
		}); err != nil {
			a.log.Error("could not send update account event", zap.Error(err))
		}
	}

	return acc, err
}
