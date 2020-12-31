package twitter

import "time"

const (
	sourceURLScheme = "https://twitter.com/api/%s.json" // TODO: fix
)

type ArchiveResponse []struct {
	Source             string  `json:"source"`
	IDStr              string  `json:"id_str"`
	Text               string  `json:"text"`
	CreatedAt          string  `json:"created_at"`
	RetweetCount       int64   `json:"retweet_count"`
	InReplyToUserIDStr *string `json:"in_reply_to_user_id_str"`
	FavoriteCount      int64   `json:"favorite_count"`
	IsRetweet          bool    `json:"is_retweet"`
}

// UserFirstTweetDate returns the date of the user first tweet.
func UserFirstTweetDate() time.Time {
	// TODO: fix.
	return time.Date(2005, time.January, 1, 0, 0, 0, 0, time.UTC)
}
