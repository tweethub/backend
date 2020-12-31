package log

import (
	"errors"

	"go.uber.org/zap"
)

// New returns a new logger.
func New(service, logLevel string) (*zap.Logger, error) {
	if service == "" {
		return nil, errors.New("missing service name")
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.MessageKey = "message"

	if err := cfg.Level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	opt := zap.Fields(zap.String("service", service))
	return cfg.Build(opt)
}
