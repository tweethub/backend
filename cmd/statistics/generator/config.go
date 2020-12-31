package generator

import (
	"errors"

	"github.com/tweethub/backend/pkg/time"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

// Config represents a generator configuration.
type Config struct {
	AllTimeStatistics  bool
	Series             int
	DefaultTimeSpans   time.Durations
	GeneratingInterval time.Duration
}

// Validate validates the generator configuration.
func (cfg *Config) Validate() error {
	var err error

	if cfg.Series <= 0 {
		err = multierr.Append(err,
			errors.New("series should be bigger than 0 in generator configuration"))
	}
	if len(cfg.DefaultTimeSpans) == 0 {
		err = multierr.Append(err,
			errors.New("missing default time spans in generator configuration"))
	}
	if cfg.GeneratingInterval <= 0 {
		err = multierr.Append(err,
			errors.New("generating interval should be bigger than 0 in generator configuration"))
	}

	return err
}

func (cfg *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("all-time-statistics", cfg.AllTimeStatistics)
	enc.AddInt("series", cfg.Series)
	enc.AddDuration("generating-interval", cfg.GeneratingInterval)

	return enc.AddArray("default-time-spans", cfg.DefaultTimeSpans)
}
