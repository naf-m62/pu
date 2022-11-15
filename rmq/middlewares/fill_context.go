package middlewares

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	wrapper "github.com/naf-m62/rabbitmq_wrapper"
	"go.uber.org/zap"

	"pu/logger"
)

func FillContext(l logger.Logger) func(m *wrapper.Middlewares) error {
	return func(m *wrapper.Middlewares) error {
		m.CtxEvent = context.WithValue(m.CtxEvent, "token", m.Delivery.MessageId)

		l = l.With(zap.String("token", m.Delivery.MessageId))
		m.CtxEvent = ctxzap.ToContext(m.CtxEvent, l)

		return m.Next()
	}
}
