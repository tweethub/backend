package app

import (
	"github.com/tweethub/backend/cmd/statistics/generator"
	"github.com/tweethub/backend/pkg/database"
	"github.com/tweethub/backend/pkg/http"
	"github.com/tweethub/backend/pkg/tweets"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

// Config defines the configuration for the service.
type Config struct {
	Generator    generator.Config
	TweetsClient tweets.Config
	RESTServer   http.Config
	Database     database.Config
}

// Validate validates the service configuration.
func (cfg *Config) Validate() error {
	return multierr.Combine(
		cfg.Generator.Validate(),
		cfg.TweetsClient.Validate(),
		cfg.RESTServer.Validate(),
		cfg.Database.Validate(),
	)
}

func (cfg *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// TODO: add the rest of the configuration.
	return enc.AddObject("generator", &cfg.Generator)
}
