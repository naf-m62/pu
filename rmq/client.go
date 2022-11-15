package rmq

import (
	"context"

	"github.com/naf-m62/rabbitmq_wrapper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pu/config"
	"pu/logger"
)

// Client only publish
type Client interface {
	ShutDown() error
	Publish(token, exchange, routingKey string, msg []byte) error
}

func NewClient(lc fx.Lifecycle, config config.Config, l logger.Logger) (r Client, err error) {
	var rCfg *wrapper.Config
	if err = config.UnmarshalKey("rabbit", &rCfg); err != nil {
		return nil, err
	}

	r, err = wrapper.New(rCfg, l, "client", nil)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
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
