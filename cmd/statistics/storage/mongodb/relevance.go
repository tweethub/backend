package mongodb

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tweethub/backend/cmd/statistics/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	keyTimeSpan  = "time_span"
	keyRelevance = "relevance"
	// ValueAllTime is the all time relevance statistics value for the time span key in the database.
	ValueAllTime = "all_time"
)

var uniqueTimeSpansIndex = mongo.IndexModel{ // nolint:gochecknoglobals
	Keys: bson.M{
		"time_span": 1,
	},
	Options: options.Index().SetName("unique_time_spans").SetUnique(true),
}

// UpdateRelevances updates the relevances into the database.
func (mng *MongoDB) UpdateRelevances(
	ctx context.Context,
	user string,
	relevances storage.Relevances,
) error {
	mng.logger.Debug("Updating relevances",
		zap.String("user", user),
		zap.Array("relevances", relevances),
	)

	if user == "" {
		return storage.ErrInvalidUser
	}

	coll := mng.db.Collection(user)

	// TODO: optimize.
	if _, err := coll.Indexes().CreateOne(ctx, uniqueTimeSpansIndex); err != nil {
		return errors.Wrap(err, "creating index")
	}

	update := bson.M{"$set": bson.M{keyRelevance: relevances}}

	upsert := true
	opts := &options.UpdateOptions{Upsert: &upsert}

	_, err := coll.UpdateOne(ctx, bson.M{}, update, opts)
	return errors.Wrap(err, "updating relevance")
}

// Relevance returns a relevance from the database.
func (mng *MongoDB) Relevance(
	ctx context.Context,
	user string,
	opts storage.RelevanceOptions,
) (*storage.Relevance, error) {
	if user == "" {
		return nil, storage.ErrInvalidUser
	}

	filter, opt := buildRelevanceQuery(opts)

	mng.logger.Debug("Getting relevance from mongodb",
		zap.String("collection", user))

	coll := mng.db.Collection(user)

	relevances := map[string]storage.Relevances{}
	err := coll.FindOne(ctx, filter, &opt).Decode(&relevances)
	if err != nil {
		return nil, errors.Wrap(err, "finding relevance")
	}

	relevance, ok := relevances[keyRelevance]
	if !ok {
		return nil, errors.New("missing relevance document")
	}

	if len(relevance) != 1 {
		return nil, errors.New("wrong relevance length")
	}
	rlv := relevance[0]

	mng.logger.Debug("Relevance from mongodb", zap.Object("relevance", rlv))

	return rlv, nil
}

func buildRelevanceQuery(opts storage.RelevanceOptions) (bson.M, options.FindOneOptions) {
	// TODO: optimize.
	var timeSpan string

	if opts.TimeSpan != nil {
		timeSpan = opts.TimeSpan.String()
	} else {
		timeSpan = ValueAllTime
	}

	filter := bson.M{
		keyRelevance + "." + keyTimeSpan: timeSpan,
	}

	opt := options.FindOneOptions{
		Projection: bson.D{
			{Key: keyID, Value: 0},
			{Key: keyRelevance, Value: bson.M{
				"$elemMatch": bson.M{
					keyTimeSpan: timeSpan,
				},
			}},
		},
	}
	return filter, opt
}
