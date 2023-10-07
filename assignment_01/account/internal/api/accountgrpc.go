package api

import (
	"context"
	"github.com/nerock/omi-test/account/internal"
	pb "github.com/nerock/omi-test/account/pb/account"
)

type AccountGrpcApi struct {
	pb.UnimplementedAccountServer

	svc internal.AccountService
}

func NewAccountGrpcApi(svc internal.AccountService) *AccountGrpcApi {
	return &AccountGrpcApi{svc: svc}
}

func (a *AccountGrpcApi) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	return nil, nil
}
