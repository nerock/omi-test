package pkg

import (
	"context"
	"errors"
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type ServicePB interface {
	RegisterPb(server *grpc.Server)
}

type GrpcServer struct {
	srv  *grpc.Server
	port int
}

func NewGrpcServer(port int, log *zap.Logger, svcs ...ServicePB) *GrpcServer {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(log)),
	)
	grpc_zap.ReplaceGrpcLoggerV2(log)

	for _, svc := range svcs {
		svc.RegisterPb(s)
	}
	reflection.Register(s)

	return &GrpcServer{
		srv:  s,
		port: port,
	}
}

func (gs *GrpcServer) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.port))
	if err != nil {
		return err
	}

	return gs.srv.Serve(lis)
}

func (gs *GrpcServer) Shutdown(ctx context.Context) error {
	closeChan := make(chan struct{})

	go func() {
		gs.srv.GracefulStop()
		closeChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("timeout closing grpc server")
	case <-closeChan:
		return nil
	}
}
