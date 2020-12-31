package v1

import (
	"time"

	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/pkg/twitter"
)

// ValidateTweetsURLValues validates the tweets URL values.
func ValidateTweetsURLValues(vls v1.TweetsURLValues) error {
	if !vls.AfterTime.IsZero() {
		if !vls.BeforeTime.IsZero() && vls.AfterTime.After(vls.BeforeTime) {
			return errInvalidBeforeTimeOrAfterTime
		}
		if vls.AfterTime.After(time.Now().UTC()) {
			return errInvalidAfterTime
		}
	}

	if !vls.BeforeTime.IsZero() && vls.BeforeTime.Before(twitter.UserFirstTweetDate()) {
		return errInvalidBeforeTime
	}
	return nil
}
