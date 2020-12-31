package app

import (
	"context"
	"time"

	"github.com/tweethub/backend/cmd/tweets/storage"
	"github.com/tweethub/backend/pkg/twitter"
	"go.uber.org/zap"
)

func (app *App) collectTweets(ctx context.Context) {
	// TODO: collect tweets from different users.
	msgStream := app.collector.CollectOverTime(ctx)

	for message := range msgStream {
		if message.Err != nil {
			app.logger.Error("Tweets message stream failed", zap.Error(message.Err))
			continue
		}

		tweets, err := archiveResponseToTweets(message.Tweets)
		if err != nil {
			app.logger.Error("Converting archives to tweets failed", zap.Error(message.Err))
			continue
		}

		app.logger.Debug("Collected tweets", zap.Int("count", len(tweets)))

		user := "" // TODO: fix
		if err := app.db.AddTweets(ctx, user, tweets); err != nil {
			app.logger.Error("Storing tweets failed", zap.Error(err))
			continue
		}
	}
}

// archiveResponseToTweets converts archive response to tweets.
func archiveResponseToTweets(resp twitter.ArchiveResponse) (storage.Tweets, error) {
	tweets := make(storage.Tweets, len(resp))

	for i, twt := range resp {
		ts, err := time.Parse(time.RubyDate, twt.CreatedAt)
		if err != nil {
			return nil, err
		}

		tweets[i] = &storage.Tweet{
			Source:      twt.Source,
			ID:          twt.IDStr,
			Text:        twt.Text,
			CreateTime:  ts.UTC(),
			Retweets:    twt.RetweetCount,
			ReplyUserID: twt.InReplyToUserIDStr,
			Favorites:   twt.FavoriteCount,
			IsRetweet:   twt.IsRetweet,
		}
	}
	return tweets, nil
}
