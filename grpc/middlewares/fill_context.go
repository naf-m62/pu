package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"pu/logger"
)

func FillContext(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		var (
			md metadata.MD
			ok bool
		)
		if md, ok = metadata.FromIncomingContext(ctx); !ok {
			l.Debug("no metadata")
			return nil, status.Error(codes.Internal, "no metadata")
		}

		tokenSlice := md.Get("x-token-id")

		var tokenStr string
		if len(tokenSlice) == 0 {
			l.Debug("no token")
			tokenStr = "default_token"
		} else {
			tokenStr = tokenSlice[0]
		}
		ctx = context.WithValue(ctx, "token", tokenStr)

		l := l.With(zap.String("token", tokenStr))
		ctx = ctxzap.ToContext(ctx, l)

		return handler(ctx, req)
	}
}
