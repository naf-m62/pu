package rmq

import (
	"context"

	"github.com/naf-m62/rabbitmq_wrapper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pu/config"
	"pu/logger"
	"pu/rmq/middlewares"
)

// Server only consume
type Server interface {
	ShutDown() error
	Consume(item *wrapper.ConsumeItem) (err error)
}

func NewServer(
	lc fx.Lifecycle, config config.Config, l logger.Logger, consumeList []wrapper.ConsumeItem,
) (r Server, err error) {
	var rCfg *wrapper.Config
	if err = config.UnmarshalKey("rabbit", &rCfg); err != nil {
		return nil, err
	}

	middlewareList := make([]func(*wrapper.Middlewares) error, 0, 3)
	middlewareList = append(middlewareList,
		middlewares.PanicRecovery(l),
		middlewares.FillContext(l),
		middlewares.ExecutionTime(l),
	)

	r, err = wrapper.New(rCfg, l, "server", middlewareList)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// list consumers
			for _, item := range consumeList {
				handler := item
				if err = r.Consume(&handler); err != nil {
					l.Error("can't start consume", zap.Error(err))
					return err
				}
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if errD := r.ShutDown(); errD != nil {
				l.Error("can't close rabbitMQ client", zap.Error(errD))
			}
			return nil
		},
	})

	return r, err
}
