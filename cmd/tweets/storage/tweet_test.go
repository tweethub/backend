package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var tweetsSorted = Tweets{{
	Source:      "1",
	ID:          "1",
	Text:        "1",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 1, 0, time.UTC),
	Retweets:    1,
	ReplyUserID: nil,
	Favorites:   1,
	IsRetweet:   false,
}, {
	Source:      "2",
	ID:          "2",
	Text:        "2",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 2, 0, time.UTC),
	Retweets:    2,
	ReplyUserID: nil,
	Favorites:   2,
	IsRetweet:   false,
}, {
	Source:      "3",
	ID:          "3",
	Text:        "3",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC),
	Retweets:    3,
	ReplyUserID: nil,
	Favorites:   3,
	IsRetweet:   false,
}}

var tweetsUnsorted = Tweets{{
	Source:      "2",
	ID:          "2",
	Text:        "2",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 2, 0, time.UTC),
	Retweets:    2,
	ReplyUserID: nil,
	Favorites:   2,
	IsRetweet:   false,
}, {
	Source:      "1",
	ID:          "1",
	Text:        "1",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 1, 0, time.UTC),
	Retweets:    1,
	ReplyUserID: nil,
	Favorites:   1,
	IsRetweet:   false,
}, {
	Source:      "3",
	ID:          "3",
	Text:        "3",
	CreateTime:  time.Date(0, 0, 0, 0, 0, 3, 0, time.UTC),
	Retweets:    3,
	ReplyUserID: nil,
	Favorites:   3,
	IsRetweet:   false,
}}

func TestTweets_SortByDate(t *testing.T) {
	tests := []struct {
		name     string
		tweets   Tweets
		expected Tweets
	}{
		{
			name:     "Test sorting tweets by date",
			tweets:   tweetsUnsorted,
			expected: tweetsSorted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tweets.SortByDate()

			require.Equal(t, tt.expected, tt.tweets)
		})
	}
}

func TestTweets_SortByFavorites(t *testing.T) {
	tests := []struct {
		name     string
		tweets   Tweets
		expected Tweets
	}{
		{
			name:     "Test sorting tweets by favorites",
			tweets:   tweetsUnsorted,
			expected: tweetsSorted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tweets.SortByFavorites()

			require.Equal(t, tt.expected, tt.tweets)
		})
	}
}

func TestTweets_SortByRetweets(t *testing.T) {
	tests := []struct {
		name     string
		tweets   Tweets
		expected Tweets
	}{
		{
			name:     "Test sorting tweets by retweets",
			tweets:   tweetsUnsorted,
			expected: tweetsSorted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tweets.SortByRetweets()

			require.Equal(t, tt.expected, tt.tweets)
		})
	}
}
