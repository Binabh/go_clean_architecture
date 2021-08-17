package bootstrap

import (
	"clean-architecture/api/controllers"
	"clean-architecture/api/middlewares"
	"clean-architecture/api/routes"
	"clean-architecture/cmd"
	"clean-architecture/infrastructure"
	"clean-architecture/lib"
	"clean-architecture/repository"
	"clean-architecture/services"
	"clean-architecture/utils"
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	routes.Module,
	services.Module,
	repository.Module,
	infrastructure.Module,
	middlewares.Module,
	cmd.Module,
	lib.Module,
	fx.Invoke(bootstrap),
)

var flushTimeout = 2 * time.Second

func bootstrap(
	lifecycle fx.Lifecycle,
	handler infrastructure.Router,
	routes routes.Routes,
	env lib.Env,
	middlewares middlewares.Middlewares,
	logger lib.Logger,
	cobracliApp cmd.RootCommands,
	database infrastructure.Database,
	migrations infrastructure.Migrations,
) {

	appStop := func(context.Context) error {
		logger.Info("Stopping Application")
		conn, _ := database.DB.DB()
		conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting hatsu cli Application")
				logger.Info("------- 🤖 clean-architecture 🤖 (CLI) -------")
				return nil
			},
			OnStop: appStop,
		})

		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting Application")
			logger.Info("-------------------------------------")
			logger.Info("------- clean-architecture 📺 -------")
			logger.Info("-------------------------------------")

			logger.Info("Migrating database schemas")
			go cobracliApp.Execute()
			//migrations.Migrate()
			go func() {
				middlewares.Setup()
				routes.Setup()
				if env.ServerPort == "" {
					handler.Run()
				} else {
					handler.Run(":" + env.ServerPort)
				}

			}()

			return sentry.Init(sentry.ClientOptions{
				Dsn:              env.SentryDSN,
				AttachStacktrace: true,
			})
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(flushTimeout)

			return nil
		},
	})
}
