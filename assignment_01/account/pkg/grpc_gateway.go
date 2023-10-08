package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nerock/omi-test/account/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type ServiceGW interface {
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type GrpcGW struct {
	conn *grpc.ClientConn
	gw   *http.Server
}

func NewGrpcGW(ctx context.Context, cfg config.ServerConfig, svcs ...ServiceGW) (*GrpcGW, error) {
	conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%d", cfg.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to grpc server: %w", err)
	}
	mux := runtime.NewServeMux()
	for _, svc := range svcs {
		if err := svc.RegisterGateway(ctx, mux, conn); err != nil {
			if errClose := conn.Close(); errClose != nil {
				err = fmt.Errorf("%w: %w", errClose, err)
			}

			return nil, err
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: mux,
	}

	return &GrpcGW{
		conn: conn,
		gw:   srv,
	}, nil
}

func (gg *GrpcGW) Serve() error {
	if err := gg.gw.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (gg *GrpcGW) Shutdown(ctx context.Context) error {
	if err := gg.gw.Shutdown(ctx); err != nil {
		return err
	}

	return gg.conn.Close()
}
