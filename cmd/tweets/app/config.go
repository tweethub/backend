package app

import (
	"github.com/tweethub/backend/pkg/database"
	"github.com/tweethub/backend/pkg/http"
	"github.com/tweethub/backend/pkg/twitter"
	"go.uber.org/multierr"
)

// Config defines the configuration for the service.
type Config struct {
	Collector  twitter.Config
	RESTServer http.Config
	Database   database.Config
}

// Validate validates the service configuration.
func (cfg *Config) Validate() error {
	return multierr.Combine(
		cfg.RESTServer.Validate(),
		cfg.Database.Validate(),
		cfg.Collector.Validate(),
	)
}
