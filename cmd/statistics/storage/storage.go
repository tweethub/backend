package storage

import (
	"context"

	"github.com/tweethub/backend/pkg/time"
	"go.uber.org/zap/zapcore"
)

// RelevanceOptions represents the criteria to get the relevance.
type RelevanceOptions struct {
	TimeSpan *time.Duration
}

// MarshalLogObject marshals relevance options.
func (opts RelevanceOptions) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("time_span", opts.TimeSpan.String())
	return nil
}

// Statistical represents storage for statistics.
type Statistical interface {
	Relevance(ctx context.Context, user string, opts RelevanceOptions) (*Relevance, error)
	UpdateRelevances(ctx context.Context, user string, relevances Relevances) error

	Close(ctx context.Context)
}
