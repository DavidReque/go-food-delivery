package grpc

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/pkg/grpc/config"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"grpcfx",
		grpcProviders,
		grpcInvokes,
	)

	grpcProviders = fx.Options(fx.Provide(
		config.ProvideConfig,
		fx.Annotate(
			NewGrpcServer,
			fx.ParamTags(``, ``),
		),
		NewGrpcClient,
	))

	grpcInvokes = fx.Options(fx.Invoke())
)

func registerHooks(
	lc fx.Lifecycle,
	grpcServer GrpcServer,
	grpcClient GrpcClient,
	logger logger.Logger,
	options *config.GrpcOptions,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				// if (ctx.Err() == nil), context not canceled or deadlined
				if err := grpcServer.RunGrpcServer(nil); err != nil {
					// do a fatal for going to OnStop process
					logger.Fatalf(
						"(GrpcServer.RunGrpcServer) error in running server: {%v}",
						err,
					)
				}
			}()

			logger.Infof(
				"%s is listening on Host:{%s} Grpc PORT: {%s}",
				options.Name,
				options.Host,
				options.Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulShutdown()
			logger.Info("server shutdown gracefully")

			if err := grpcClient.Close(); err != nil {
				logger.Errorf("error in closing grpc client: {%v}", err)
			} else {
				logger.Info("grpc-client closed gracefully")
			}

			return nil
		},
	})
}
