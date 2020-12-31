package tweets

import (
	"errors"

	"go.uber.org/multierr"
)

// Config represents tweets client configuration.
type Config struct {
	Source     string
	APIVersion string
	tweetsURL  string
}

// Validate validates the tweets client configuration.
func (cfg Config) Validate() error {
	var err error

	if cfg.Source == "" {
		err = multierr.Append(err, errors.New("missing source in tweets client configuration"))
	}
	if cfg.APIVersion == "" {
		err = multierr.Append(err, errors.New("missing API version in tweets client configuration"))
	}

	return err
}
