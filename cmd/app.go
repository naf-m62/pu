package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"pu/cmd/database"
	grpcprivate "pu/cmd/handlers/grpc_handlers/private"
	grpcpublic "pu/cmd/handlers/grpc_handlers/public"
	httphandlers "pu/cmd/handlers/http_handlers"
	rmqhandlers "pu/cmd/handlers/rmq_handlers"
	"pu/cmd/processor/user"
	"pu/cmd/publisher"
	"pu/config"
	"pu/grpc"
	privateproto "pu/grpc/proto/private"
	publicproto "pu/grpc/proto/public"
	"pu/http"
	"pu/logger"
	"pu/postgres"
	"pu/rmq"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "run app",
	Run:   runApp,
}

func runApp(_ *cobra.Command, _ []string) {
	app := fx.New(
		fx.Provide(
			fx.Annotate(getConfigPath, fx.ResultTags(`name:"configPath"`)),
			fx.Annotate(config.New, fx.ParamTags(`name:"configPath"`)),
			logger.New,

			postgres.New,
			database.NewTx,
			database.NewUser,

			rmq.NewClient,
			publisher.NewPublisher,

			user.NewProcessor,

			rmqhandlers.NewHandler,
			rmqhandlers.RegisterHandlerList,

			httphandlers.NewHandler,
			httphandlers.RegisterHandlerList,

			fx.Annotate(grpc.NewGrpcPublicServer, fx.ResultTags(`name:"grpc_public_def"`)),
			fx.Annotate(grpcpublic.NewGrpcPublicServer, fx.ResultTags(`name:"grpc_public"`)),

			fx.Annotate(grpc.NewGrpcPrivateServer, fx.ResultTags(`name:"grpc_private_def"`)),
			fx.Annotate(grpcprivate.NewGrpcPrivateServer, fx.ResultTags(`name:"grpc_private"`)),
		),
		// grpc/http clients

		// watchers

		// rmq servers
		fx.Invoke(rmq.NewServer),
		// grpc/http servers
		fx.Invoke(http.NewHTTPServer),

		fx.Invoke(fx.Annotate(privateproto.RegisterGRPCPrivateServer, fx.ParamTags(`name:"grpc_private_def"`, `name:"grpc_private"`))),
		fx.Invoke(fx.Annotate(publicproto.RegisterGRPCPublicServer, fx.ParamTags(`name:"grpc_public_def"`, `name:"grpc_public"`))),
	)

	startCtx, startCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	<-app.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
