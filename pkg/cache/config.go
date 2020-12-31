package cache

import "errors"

// Config represents cache configuration.
type Config struct {
	Address  string
	Password string
	DB       int
}

// Validate validates the cache configuration.
func (cfg Config) Validate() error {
	if cfg.Address == "" {
		return errors.New("missing address in cache configuration")
	}
	return nil
}
