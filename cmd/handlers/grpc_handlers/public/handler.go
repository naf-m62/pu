package public

import (
	"context"

	"pu/cmd/handlers/grpc_handlers/errors"

	"pu/cmd/processor/user"
	"pu/cmd/processor/user/domain"
	publicproto "pu/grpc/proto/public"
)

type grpcServer struct {
	publicproto.UnimplementedGRPCPublicServer

	userProcessor *user.Processor
}

func NewGrpcPublicServer(userProcessor *user.Processor) publicproto.GRPCPublicServer {
	return &grpcServer{userProcessor: userProcessor}
}

func (s *grpcServer) CreateUser(
	ctx context.Context, req *publicproto.CreateUserRequest,
) (_ *publicproto.User, err error) {
	var up *domain.User
	if up, err = s.userProcessor.Create(ctx, &domain.User{Name: req.Name, Email: req.Email}, req.Password); err != nil {
		return nil, errors.GetError(err)
	}
	return convertUserToProto(up), nil
}

func (s *grpcServer) AuthUser(ctx context.Context, req *publicproto.AuthUserRequest) (_ *publicproto.User, err error) {
	var up *domain.User
	if up, err = s.userProcessor.Auth(ctx, req.Email, req.EncodedPassword); err != nil {
		return nil, errors.GetError(err)
	}
	return convertUserToProto(up), nil
}

func (s *grpcServer) GetUser(ctx context.Context, req *publicproto.GetUserRequest) (res *publicproto.User, err error) {
	var up *domain.User
	if up, err = s.userProcessor.Get(ctx, req.Id); err != nil {
		return nil, errors.GetError(err)
	}
	return convertUserToProto(up), nil
}

func convertUserToProto(u *domain.User) *publicproto.User {
	return &publicproto.User{
		Id:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Points: u.Points,
	}
}
