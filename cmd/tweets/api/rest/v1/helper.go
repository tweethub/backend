package v1

import (
	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/cmd/tweets/storage"
)

func GenerateTweetsOptions(vls v1.TweetsURLValues) storage.TweetsOptions {
	opts := storage.TweetsOptions{}

	if !vls.AfterTime.IsZero() {
		opts.AfterTime = &vls.AfterTime
	}
	if !vls.BeforeTime.IsZero() {
		opts.BeforeTime = &vls.BeforeTime
	}
	return opts
}
