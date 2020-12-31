package app

import (
	"context"

	rest "github.com/tweethub/backend/cmd/tweets/api/rest/v1"
	"github.com/tweethub/backend/cmd/tweets/storage"
	"github.com/tweethub/backend/cmd/tweets/storage/mongodb"
	"github.com/tweethub/backend/pkg/log"
	"github.com/tweethub/backend/pkg/service"
	"github.com/tweethub/backend/pkg/twitter"
	"go.uber.org/zap"
)

const serviceName = "tweets"

// App represents an application.
type App struct {
	logger    *zap.Logger
	collector *twitter.Collector
	restSrv   *rest.Server
	db        storage.Twitter
}

// Init initializes an application.
func Init(ctx context.Context, configPath, logLevel string) *App {
	logger, err := log.New(serviceName, logLevel)
	if err != nil {
		panic(err)
	}

	logger.Info(service.Init)

	cfg := Config{}
	service.ReadConfig(configPath, &cfg, logger)

	db := mongodb.Init(ctx, cfg.Database, logger)
	restSrv := rest.NewServer(db, cfg.RESTServer, logger)
	cltr := twitter.NewCollector(cfg.Collector, logger)

	return &App{
		logger:    logger,
		restSrv:   restSrv,
		db:        db,
		collector: cltr,
	}
}

// Start starts the application.
func (app *App) Start(ctx context.Context) {
	app.logger.Info(service.Running)

	go app.restSrv.Start(ctx)
	go app.collectTweets(ctx)
}

// Stop stops the application.
func (app *App) Stop(ctx context.Context) {
	app.logger.Info(service.ShuttingDown)

	app.db.Close(ctx)
}
