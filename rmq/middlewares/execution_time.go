package middlewares

import (
	"strconv"
	"time"

	"github.com/naf-m62/rabbitmq_wrapper"

	"pu/logger"
)

func ExecutionTime(l logger.Logger) func(m *wrapper.Middlewares) error {
	return func(m *wrapper.Middlewares) error {
		ts := time.Now()
		err := m.Next()
		te := time.Now()
		l.Debug("executed: " + strconv.FormatInt(te.UnixMilli()-ts.UnixMilli(), 10) + " milliseconds")
		return err
	}
}
