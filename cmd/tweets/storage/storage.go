package storage

import (
	"context"
	"time"
)

// TweetsOptions represents the criteria to get the tweets.
type TweetsOptions struct {
	AfterTime   *time.Time
	BeforeTime  *time.Time
	IsRetweet   *bool
	ReplyUserID *string
}

// Twitter represents storage for tweets.
type Twitter interface {
	AddTweets(ctx context.Context, user string, tweets Tweets) error
	Tweets(ctx context.Context, user string, opts TweetsOptions) (Tweets, error)

	Close(ctx context.Context)
}
