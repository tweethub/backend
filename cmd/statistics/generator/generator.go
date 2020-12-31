package generator

import (
	"context"
	gotime "time"

	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"github.com/tweethub/backend/pkg/time"
	"github.com/tweethub/backend/pkg/tweets"
	"github.com/tweethub/backend/pkg/twitter"
	"go.uber.org/zap"
)

// Generator represents a statistics generator.
type Generator struct {
	config     Config
	logger     *zap.Logger
	twtsClient *tweets.Client
	db         storage.Statistical
}

// New returns new generator.
func New(db storage.Statistical, twtsClient *tweets.Client, config Config, logger *zap.Logger) *Generator {
	config.DefaultTimeSpans.Sort()

	return &Generator{
		config:     config,
		logger:     logger,
		twtsClient: twtsClient,
		db:         db,
	}
}

// Start starts the generator.
func (gen *Generator) Start(ctx context.Context) {
	gen.logger.Info("Generating statistics")

	ticker := gotime.NewTicker(gen.config.GeneratingInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go gen.GenerateRelevances(ctx)
		}
	}
}

func (gen *Generator) tweets(ctx context.Context) (v1.TweetsResponse, error) {
	firstTweetDate := twitter.UserFirstTweetDate()

	timeNow := time.NowUTC()
	largestTimeSpan := gen.config.DefaultTimeSpans[len(gen.config.DefaultTimeSpans)-1]
	largestTimeSpanStartDate := timeNow.Add(-largestTimeSpan)

	var urlValues v1.TweetsURLValues
	if !gen.config.AllTimeStatistics && largestTimeSpanStartDate.After(firstTweetDate) {
		// Get only the tweets after the start of the largest time span.
		urlValues.AfterTime = largestTimeSpanStartDate
	}

	user := "" // TODO: fix
	return gen.twtsClient.GetTweets(ctx, user, urlValues)
}
