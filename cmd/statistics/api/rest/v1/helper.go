package v1

import (
	statistics "github.com/tweethub/backend/api/services/statistics/v1"
	tweets "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/pkg/time"
	"github.com/tweethub/backend/pkg/twitter"
)

// GenerateTweetsURLValues generates tweets URL values from relevance URL values.
func GenerateTweetsURLValues(vls statistics.RelevanceURLValues) tweets.TweetsURLValues {
	values := tweets.TweetsURLValues{}

	if vls.BeforeTime.IsZero() && vls.AfterTime.IsZero() {
		return values
	}
	if vls.BeforeTime.IsZero() {
		values.BeforeTime = time.NowUTC()
	}
	if vls.AfterTime.IsZero() {
		values.AfterTime = twitter.UserFirstTweetDate()
	}
	return values
}

// GenerateTimeFrame generates time frame from the URL values.
func GenerateTimeFrame(vls statistics.RelevanceURLValues) time.Frame {
	timeFrame := time.Frame{
		StartTime: twitter.UserFirstTweetDate(),
		EndTime:   time.NowUTC(),
	}

	if !vls.BeforeTime.IsZero() {
		timeFrame.EndTime = vls.BeforeTime
	}
	if !vls.AfterTime.IsZero() {
		timeFrame.StartTime = vls.AfterTime
	}
	return timeFrame
}
