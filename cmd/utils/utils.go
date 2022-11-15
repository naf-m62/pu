package utils

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"pu/logger"
)

func GetTokenFromContext(ctx context.Context) string {
	tokenAny := ctx.Value("token")

	token := "default_token"
	if tmp, ok := tokenAny.(string); ok {
		token = tmp
	}
	return token
}

func GetLoggerFromContext(ctx context.Context) logger.Logger {
	return ctxzap.Extract(ctx)
}
