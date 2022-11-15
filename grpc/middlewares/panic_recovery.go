package middlewares

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pu/logger"
)

func PanicRecovery(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				l.Error("panic happened: " + fmt.Sprint(rec))
				err = status.Error(codes.Internal, "")
			}
		}()
		return handler(ctx, req)
	}
}
