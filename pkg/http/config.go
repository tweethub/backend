package http

import (
	"errors"

	"go.uber.org/multierr"
)

// Config represents server configuration.
type Config struct {
	Host     string
	Port     string
	CertFile string
	KeyFile  string
}

// Validate validates the server configuration.
func (cfg Config) Validate() error {
	var err error

	if cfg.Host == "" {
		err = multierr.Append(err, errors.New("missing host in server configuration"))
	}
	if cfg.Port == "" {
		err = multierr.Append(err, errors.New("missing port in server configuration"))
	}

	return err
}
