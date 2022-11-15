package middlewares

import (
	"context"

	"google.golang.org/grpc"

	"pu/logger"
)

// todo
func Auth(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		l.Debug("user auth")
		return handler(ctx, req)
	}
}
