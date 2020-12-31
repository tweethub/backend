package generator

import (
	"context"

	"github.com/pkg/errors"
	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"github.com/tweethub/backend/cmd/statistics/storage/mongodb"
	"github.com/tweethub/backend/pkg/time"
	"github.com/tweethub/backend/pkg/twitter"
	"go.uber.org/zap"
)

// GenerateRelevances gets all the tweets, generates relevance statistics for all time spans and
// updates the relevances in the database.
func (gen *Generator) GenerateRelevances(ctx context.Context) {
	gen.logger.Debug("Generating relevance")

	tweets, err := gen.tweets(ctx)
	if err != nil {
		gen.logger.Error("Getting tweets failed", zap.Error(err))
		return
	}

	relevances, err := gen.CalculateRelevances(tweets)
	if err != nil {
		gen.logger.Error("Calculating relevances failed", zap.Error(err))
		return
	}

	user := "" // TODO: fix.

	if err := gen.db.UpdateRelevances(ctx, user, relevances); err != nil {
		gen.logger.Error("Updating relevances failed",
			zap.Error(err),
			zap.String("user", user),
			zap.Array("relevances", relevances),
		)
		return
	}
}

// CalculateRelevances calculates the relevances for different time spans.
func (gen *Generator) CalculateRelevances(tweets v1.TweetsResponse) (storage.Relevances, error) {
	if len(tweets) == 0 {
		return nil, errors.New("missing tweets")
	}

	var relevances storage.Relevances

	timeFrame := time.Frame{
		EndTime: time.NowUTC(),
	}

	tweets.SortByDate()

	for _, timeSpan := range gen.config.DefaultTimeSpans {
		timeFrame.StartTime = timeFrame.EndTime.Add(-timeSpan)
		relevance, err := CalculateRelevance(tweets, timeFrame, gen.config.Series)
		if err != nil {
			return nil, err
		}
		relevances = append(relevances, relevance)
	}

	if gen.config.AllTimeStatistics {
		timeFrame.StartTime = twitter.UserFirstTweetDate() // TODO: add support for more users.
		relevance, err := CalculateRelevance(tweets, timeFrame, gen.config.Series)
		if err != nil {
			return nil, err
		}
		relevance.TimeSpan = mongodb.ValueAllTime // HACK: move
		relevances = append(relevances, relevance)
	}
	return relevances, nil
}

// CalculateRelevance calculates relevance.
func CalculateRelevance(
	tweets v1.TweetsResponse,
	timeFrame time.Frame,
	series int,
) (*storage.Relevance, error) {
	if timeFrame.StartTime.After(timeFrame.EndTime) {
		return nil, errors.New("start date is after end date")
	}
	if timeFrame.EndTime.Before(timeFrame.StartTime) {
		return nil, errors.New("end date is before start date")
	}
	if timeFrame.StartTime.After(time.NowUTC()) {
		return nil, errors.New("start date is in the future")
	}
	if timeFrame.StartTime.Equal(timeFrame.EndTime) {
		return nil, errors.New("the start date is the same as the end date")
	}

	relevanceCandles := calculateRelevanceCandles(tweets, timeFrame, series)

	return &storage.Relevance{
		TimeSpan: timeFrame.ToDurationString(),
		Candles:  relevanceCandles,
	}, nil
}

// calculateRelevanceCandles calculates the relevance candles.
func calculateRelevanceCandles(
	tweets v1.TweetsResponse,
	timeFrame time.Frame,
	series int,
) []*storage.RelevanceCandle {
	timeFrameDuration := timeFrame.Duration()
	durationPerTimeFrame := timeFrameDuration / time.Duration(series)

	timeFrames := time.GenerateFrames(timeFrame.StartTime, durationPerTimeFrame, series)
	relevanceCandles := make([]*storage.RelevanceCandle, series)

	for i, timeFrame := range timeFrames {
		relevanceCandle := calculateRelevanceCandle(tweets, *timeFrame)
		relevanceCandles[i] = relevanceCandle
	}
	return relevanceCandles
}

// calculateRelevanceCandle calculates relevance candle.
func calculateRelevanceCandle(tweets v1.TweetsResponse, timeFrame time.Frame) *storage.RelevanceCandle {
	relevanceCandle := storage.NewRelevanceCandle()
	relevanceCandle.TimeFrame = timeFrame.ToUnixString()

	twts := tweetsInTimeFrame(tweets, timeFrame)

	tweetsLen := len(twts)
	if tweetsLen == 0 {
		return relevanceCandle
	}
	relevanceCandle.TweetsCount = int64(tweetsLen)

	calculateFavorites(relevanceCandle, twts)
	calculateRetweets(relevanceCandle, twts)

	return relevanceCandle
}

func calculateRetweets(relevanceCandle *storage.RelevanceCandle, twts v1.TweetsResponse) { // nolint:dupl
	tweetsLen := len(twts)

	twts.SortByRetweets()

	relevanceCandle.Retweets.High = twts[tweetsLen-1].Retweets
	relevanceCandle.Retweets.Low = twts[0].Retweets

	if tweetsLen%2 == 0 {
		ret1 := twts[tweetsLen/2].Retweets
		ret2 := twts[tweetsLen/2-1].Retweets
		relevanceCandle.Retweets.Median = float64(ret1+ret2) / 2 // nolint:gomnd
	} else {
		relevanceCandle.Retweets.Median = float64(twts[tweetsLen/2].Retweets)
	}
}

func calculateFavorites(relevanceCandle *storage.RelevanceCandle, twts v1.TweetsResponse) { // nolint:dupl
	tweetsLen := len(twts)

	twts.SortByFavorites()

	relevanceCandle.Favorites.High = twts[tweetsLen-1].Favorites
	relevanceCandle.Favorites.Low = twts[0].Favorites

	if tweetsLen%2 == 0 {
		fav1 := twts[tweetsLen/2].Favorites
		fav2 := twts[tweetsLen/2-1].Favorites
		relevanceCandle.Favorites.Median = float64(fav1+fav2) / 2 // nolint:gomnd
	} else {
		relevanceCandle.Favorites.Median = float64(twts[tweetsLen/2].Favorites)
	}
}

func tweetsInTimeFrame(tweets v1.TweetsResponse, timeFrame time.Frame) v1.TweetsResponse {
	var twts v1.TweetsResponse

	for _, tweet := range tweets {
		if tweet.CreateTime.Before(timeFrame.StartTime) || tweet.CreateTime.After(timeFrame.EndTime) {
			continue
		}
		twts = append(twts, tweet)
	}
	return twts
}
