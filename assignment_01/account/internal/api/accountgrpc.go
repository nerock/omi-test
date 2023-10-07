package api

import (
	"context"
	"errors"
	"github.com/nerock/omi-test/account/internal"
	pb "github.com/nerock/omi-test/account/pb/account"
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
