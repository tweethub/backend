package main

import (
	"context"
	"flag"

	"github.com/tweethub/backend/cmd/tweets/app"
)

func main() {
	ctx := context.Background()

	var configPath string
	var logLevel string
	flag.StringVar(&configPath, "p", "prod.toml", "configuration file path")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.Parse()

	application := app.Init(ctx, configPath, logLevel)
	application.Start(ctx)
	defer application.Stop(ctx)

	<-ctx.Done()
}
