package v1

import (
	v1 "github.com/tweethub/backend/api/services/statistics/v1"
	"github.com/tweethub/backend/pkg/time"
	"github.com/tweethub/backend/pkg/twitter"
)

// ValidateRelevanceURLValues validates the relevance URL values.
func ValidateRelevanceURLValues(vls v1.RelevanceURLValues) error {
	if vls.Series <= 0 {
		return errInvalidSeries
	}
	if !vls.BeforeTime.IsZero() && vls.BeforeTime.Before(twitter.UserFirstTweetDate()) {
		return errInvalidBeforeTime
	}
	if !vls.BeforeTime.IsZero() && vls.AfterTime.After(vls.BeforeTime) {
		return errInvalidBeforeTimeOrAfterTime
	}
	if vls.AfterTime.After(time.NowUTC()) {
		return errInvalidAfterTime
	}
	return nil
}
