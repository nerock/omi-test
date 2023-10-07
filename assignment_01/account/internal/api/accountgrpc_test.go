package api

import (
	"context"
	"errors"
	"github.com/nerock/omi-test/account/internal"
	pb "github.com/nerock/omi-test/account/pb/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockAccountService struct {
	updateAccountFunc func(ctx context.Context, id string, name string) (internal.Account, error)
}

func (m *mockAccountService) UpdateAccount(ctx context.Context, id, name string) (internal.Account, error) {
	return m.updateAccountFunc(ctx, id, name)
}

func TestAccountGrpcApi_UpdateAccount(t *testing.T) {
	Convey("UpdateAccount", t, func() {
		ctx := context.Background()
		svc := &mockAccountService{}
		api := NewAccountGrpcApi(svc)

		Convey("when the request is invalid", func() {
			req := &pb.UpdateAccountRequest{}

			Convey("the api should return a nil response and a validation error", func() {
				res, err := api.UpdateAccount(ctx, req)

				So(res, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(status.Code(err), ShouldEqual, codes.InvalidArgument)
			})
		})

		Convey("when the request is valid", func() {
			req := &pb.UpdateAccountRequest{
				Id:   "id",
				Name: "name",
			}

			Convey("and the service returns a validation error", func() {
				svc.updateAccountFunc = func(inCtx context.Context, inID string, inName string) (internal.Account, error) {
					So(inCtx, ShouldResemble, ctx)
					So(inID, ShouldEqual, req.Id)
					So(inName, ShouldEqual, req.Name)

					return internal.Account{}, internal.ErrInvalidAccount
				}

				Convey("the api should return a nil response and a validation error", func() {
					res, err := api.UpdateAccount(ctx, req)

					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(status.Code(err), ShouldEqual, codes.InvalidArgument)
				})
			})

			Convey("and the service returns an unknown error", func() {
				svc.updateAccountFunc = func(inCtx context.Context, inID string, inName string) (internal.Account, error) {
					So(inCtx, ShouldResemble, ctx)
					So(inID, ShouldEqual, req.Id)
					So(inName, ShouldEqual, req.Name)

					return internal.Account{}, errors.New("something happened")
				}

				Convey("the api should return a nil response and an internal error", func() {
					res, err := api.UpdateAccount(ctx, req)

					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(status.Code(err), ShouldEqual, codes.Internal)
				})
			})

			Convey("and the service returns the updated account and no error", func() {
				svc.updateAccountFunc = func(inCtx context.Context, inID string, inName string) (internal.Account, error) {
					So(inCtx, ShouldResemble, ctx)
					So(inID, ShouldEqual, req.Id)
					So(inName, ShouldEqual, req.Name)

					return internal.Account{ID: req.Id, Name: req.Name}, nil
				}

				Convey("the api should return a successful response and no error", func() {
					res, err := api.UpdateAccount(ctx, req)

					So(res, ShouldResemble, &pb.UpdateAccountResponse{Success: true})
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
