package private

import (
	"context"

	"pu/cmd/handlers/grpc_handlers/errors"
	"pu/cmd/processor/user"
	"pu/cmd/processor/user/domain"
	privateproto "pu/grpc/proto/private"
)

type grpcServer struct {
	privateproto.UnimplementedGRPCPrivateServer

	userProcessor *user.Processor
}

func NewGrpcPrivateServer(processor *user.Processor) privateproto.GRPCPrivateServer {
	return &grpcServer{userProcessor: processor}
}

func (s *grpcServer) GetUser(ctx context.Context, req *privateproto.GetUserRequest) (_ *privateproto.User, err error) {
	var up *domain.User
	if up, err = s.userProcessor.Get(ctx, req.Id); err != nil {
		return nil, errors.GetError(err)
	}
	return convertUserToProto(up), nil
}

func convertUserToProto(u *domain.User) *privateproto.User {
	return &privateproto.User{
		Id:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Points: u.Points,
	}
}
