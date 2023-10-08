package api

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nerock/omi-test/account/internal"
	pb "github.com/nerock/omi-test/account/pb/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountGrpcApi struct {
	pb.UnimplementedAccountServer

	svc internal.AccountService
}

func NewAccountGrpcApi(svc internal.AccountService) *AccountGrpcApi {
	return &AccountGrpcApi{svc: svc}
}

func (a *AccountGrpcApi) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	if req.GetId() == "" || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid account parameters")
	}

	if _, err := a.svc.UpdateAccount(ctx, req.GetId(), req.GetName()); errors.Is(err, internal.ErrInvalidAccount) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, "could not update account")
	}

	return &pb.UpdateAccountResponse{Success: true}, nil
}

func (a *AccountGrpcApi) RegisterPb(s *grpc.Server) {
	pb.RegisterAccountServer(s, a)
}

func (a *AccountGrpcApi) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	if err := pb.RegisterAccountHandler(ctx, mux, conn); err != nil {
		return err
	}

	return nil
}
