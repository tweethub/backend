package storage

import (
	"sort"
	"time"

	"go.uber.org/zap/zapcore"
)

// Tweet represents a tweet.
type Tweet struct {
	Source      string    `json:"source" bson:"source" example:"Twitter Web Client"`
	ID          string    `json:"id" bson:"id" example:"4629116949"`
	Text        string    `json:"text" bson:"text" example:"- Read what Peter has to say about ..."`
	CreateTime  time.Time `json:"create_time" bson:"create_time" example:"2010-03-03T14:37:38.000Z"`
	Retweets    int64     `json:"retweets" bson:"retweets" example:"1"`
	ReplyUserID *string   `json:"reply_user_id" bson:"reply_user_id" example:"996852" extensions:"x-nullable"`
	Favorites   int64     `json:"favorites" bson:"favorites" example:"4"`
	IsRetweet   bool      `json:"is_retweet" bson:"is_retweet" example:"false"`
}

// MarshalLogObject marshals tweet.
func (twt *Tweet) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("source", twt.Source)
	enc.AddString("id", twt.ID)
	enc.AddString("text", twt.Text)
	enc.AddTime("create_time", twt.CreateTime)
	enc.AddInt64("retweets", twt.Retweets)
	enc.AddInt64("favorites", twt.Favorites)
	enc.AddBool("is_retweet", twt.IsRetweet)

	if twt.ReplyUserID != nil {
		enc.AddString("reply_user_id", *twt.ReplyUserID)
	} else {
		enc.AddString("reply_user_id", "")
	}
	return nil
}

// Tweets represents a list of tweets.
type Tweets []*Tweet

func (twts Tweets) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, twt := range twts {
		err := enc.AppendObject(twt)
		if err != nil {
			return err
		}
	}
	return nil
}

// SortByDate sorts tweets by date.
func (twts Tweets) SortByDate() {
	sort.Slice(twts, func(i, j int) bool {
		return twts[i].CreateTime.Before(twts[j].CreateTime)
	})
}

// SortByFavorites sorts tweets by favorites count.
func (twts Tweets) SortByFavorites() {
	sort.Slice(twts, func(i, j int) bool {
		return twts[i].Favorites < twts[j].Favorites
	})
}

// SortByRetweets sorts tweets by retweets count.
func (twts Tweets) SortByRetweets() {
	sort.Slice(twts, func(i, j int) bool {
		return twts[i].Retweets < twts[j].Retweets
	})
}
