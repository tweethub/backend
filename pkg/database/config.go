package database

import (
	"errors"

	"go.uber.org/multierr"
)

// Config represents a database configuration.
type Config struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

// Validate checks if the database configuration is correct.
func (cfg *Config) Validate() error {
	var err error

	if cfg.Name == "" {
		err = multierr.Append(err, errors.New("missing database name in the database configuration"))
	}
	if cfg.User == "" {
		err = multierr.Append(err, errors.New("missing user in the database configuration"))
	}
	if cfg.Host == "" {
		err = multierr.Append(err, errors.New("missing host in the database configuration"))
	}
	if cfg.Port == "" {
		err = multierr.Append(err, errors.New("missing port in the database configuration"))
	}
	if cfg.Password == "" {
		err = multierr.Append(err, errors.New("missing password in the database configuration"))
	}

	return err
}
