package middlewares

import (
	"fmt"

	"github.com/naf-m62/rabbitmq_wrapper"

	"pu/logger"
)

func PanicRecovery(l logger.Logger) func(m *wrapper.Middlewares) error {
	return func(m *wrapper.Middlewares) error {
		defer func() {
			if rec := recover(); rec != nil {
				l.Error("panic happened:" + fmt.Sprint(rec))
			}
		}()
		return m.Next()
	}
}
