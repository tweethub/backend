package service

import (
	"github.com/hapci/go-config"
	"go.uber.org/zap"
)

type Config interface {
	Validate() error
}

// ReadConfig parses the configuration file and validates the configuration values.
func ReadConfig(configSource string, cfg Config, logger *zap.Logger) {
	logger.Info(ReadingConfig)

	if err := config.UnmarshalFromFile(configSource, &cfg); err != nil {
		logger.Fatal(ReadingConfigFailed, zap.Error(err))
	}

	if err := cfg.Validate(); err != nil {
		logger.Fatal(ValidatingConfigFailed, zap.Error(err))
	}
}
