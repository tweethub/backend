package mongodb

import (
	"context"
	goerrors "errors"

	"github.com/pkg/errors"
	"github.com/tweethub/backend/cmd/tweets/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	keyCreateTime  = "create_time"
	keyIsRetweet   = "is_retweet"
	keyReplyUserID = "reply_user_id"
)

var uniqueIDsIndex = mongo.IndexModel{ // nolint:gochecknoglobals
	Keys: bson.M{
		"id": 1,
	},
	Options: options.Index().SetName("unique_ids").SetUnique(true),
}

// AddTweets stores tweets to the database.
func (mng *MongoDB) AddTweets(ctx context.Context, user string, tweets storage.Tweets) error {
	if user == "" {
		return storage.ErrInvalidUser
	}

	coll := mng.db.Collection(user)

	// TODO: optimize.
	if _, err := coll.Indexes().CreateOne(ctx, uniqueIDsIndex); err != nil {
		return err
	}

	twts := make([]interface{}, len(tweets))
	for i, t := range tweets {
		twts[i] = t
	}

	opts := &options.InsertManyOptions{
		Ordered: new(bool),
	}

	_, err := coll.InsertMany(ctx, twts, opts)
	if !goerrors.As(err, &mongo.BulkWriteException{}) { // Ignores duplicate entries error.
		return err
	}
	return nil
}

// Tweets returns tweets from the database.
func (mng *MongoDB) Tweets(ctx context.Context, user string, opts storage.TweetsOptions) (storage.Tweets, error) {
	if user == "" {
		return nil, storage.ErrInvalidUser
	}

	coll := mng.db.Collection(user)

	filter := buildTweetsFilter(opts)

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "finding tweets")
	}

	var tweets storage.Tweets
	err = cur.All(ctx, &tweets)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

// buildTweetsFilter builds a mongodb filter for tweets data.
func buildTweetsFilter(opts storage.TweetsOptions) interface{} {
	filter := bson.M{}

	createTimeFilter := bson.M{}
	if opts.AfterTime != nil {
		createTimeFilter["$gte"] = opts.AfterTime.UTC()
	}
	if opts.BeforeTime != nil {
		createTimeFilter["$lte"] = opts.BeforeTime.UTC()
	}
	if len(createTimeFilter) != 0 {
		filter[keyCreateTime] = createTimeFilter
	}

	if opts.IsRetweet != nil {
		filter[keyIsRetweet] = bson.M{
			"$eq": *opts.IsRetweet,
		}
	}

	if opts.ReplyUserID != nil {
		if *opts.ReplyUserID == "" {
			filter[keyReplyUserID] = nil
		} else {
			filter[keyReplyUserID] = bson.M{
				"$eq": *opts.ReplyUserID,
			}
		}
	}
	return filter
}
