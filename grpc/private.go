package grpc

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"pu/config"
	"pu/grpc/middlewares"
	"pu/logger"
)

func NewGrpcPrivateServer(
	lc fx.Lifecycle, config config.Config, log logger.Logger,
) (_ grpc.ServiceRegistrar, err error) {
	gConf := &grpcConfig{}
	if err = config.UnmarshalKey("grpc.private", &gConf); err != nil {
		return nil, err
	}

	middlewareList := []grpc.UnaryServerInterceptor{
		middlewares.PanicRecovery(log),
		middlewares.FillContext(log),
	}
	return grpcServer(lc, gConf, log, "private", middlewareList)
}
