package app

import (
	"context"

	rest "github.com/tweethub/backend/cmd/statistics/api/rest/v1"
	"github.com/tweethub/backend/cmd/statistics/generator"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"github.com/tweethub/backend/cmd/statistics/storage/mongodb"
	"github.com/tweethub/backend/pkg/log"
	"github.com/tweethub/backend/pkg/service"
	"github.com/tweethub/backend/pkg/tweets"
	"go.uber.org/zap"
)

const serviceName = "statistics"

// App represents an application.
type App struct {
	logger    *zap.Logger
	generator *generator.Generator
	restSrv   *rest.Server
	db        storage.Statistical
}

// Init initializes the application.
func Init(ctx context.Context, configPath, logLevel string) *App {
	logger, err := log.New(serviceName, logLevel)
	if err != nil {
		panic(err)
	}

	logger.Info(service.Init)

	cfg := Config{}
	service.ReadConfig(configPath, &cfg, logger)

	db := mongodb.Init(ctx, cfg.Database, logger)
	twtsClient := tweets.NewClient(cfg.TweetsClient, logger)
	restSrv := rest.NewServer(db, twtsClient, cfg.RESTServer, logger)
	gen := generator.New(db, twtsClient, cfg.Generator, logger)

	app := &App{
		logger:    logger,
		generator: gen,
		restSrv:   restSrv,
		db:        db,
	}
	return app
}

// Start starts the application.
func (app *App) Start(ctx context.Context) {
	app.logger.Info(service.Running)

	go app.restSrv.Start(ctx)
	go app.generator.Start(ctx)
}

// Stop stops the application.
func (app *App) Stop(ctx context.Context) {
	app.logger.Info(service.ShuttingDown)

	app.db.Close(ctx)
}
