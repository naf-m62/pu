package middlewares

import (
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func FillContext(h http.HandlerFunc, l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-token-id")
		if token == "" {
			token = "default_token"
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", token)

		l := l.With(zap.String("token", token))
		ctx = ctxzap.ToContext(ctx, l)

		r = r.WithContext(ctx)

		h(w, r)
	}
}
