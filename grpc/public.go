package grpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"pu/config"
	"pu/grpc/middlewares"
	"pu/logger"
)

type grpcConfig struct {
	Addr string `mapstructure:"addr"`
}

func NewGrpcPublicServer(
	lc fx.Lifecycle, config config.Config, log logger.Logger,
) (_ grpc.ServiceRegistrar, err error) {
	gConf := &grpcConfig{}
	if err = config.UnmarshalKey("grpc.public", &gConf); err != nil {
		return nil, err
	}

	middlewareList := []grpc.UnaryServerInterceptor{
		middlewares.PanicRecovery(log),
		middlewares.FillContext(log),
		middlewares.Auth(log),
	}

	return grpcServer(lc, gConf, log, "public", middlewareList)
}

func grpcServer(
	lc fx.Lifecycle, gConf *grpcConfig, log logger.Logger, name string, middlewareList []grpc.UnaryServerInterceptor,
) (_ grpc.ServiceRegistrar, err error) {
	gs := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(middlewareList...)))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var listener net.Listener
			if listener, err = net.Listen("tcp", gConf.Addr); err != nil {
				log.Error("error net.Listen", zap.Error(err))
				return err
			}
			log.Info("Start GRPC " + name + " server on " + gConf.Addr)
			go func() {
				if err = gs.Serve(listener); err != nil {
					log.Error("grpc server can't listen and serve", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			gs.GracefulStop()
			return nil
		},
	},
	)
	return func() grpc.ServiceRegistrar { return gs }(), nil
}
