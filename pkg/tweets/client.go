package tweets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	v1 "github.com/tweethub/backend/api/services/tweets/v1"
	"github.com/tweethub/backend/pkg/json"
	"go.uber.org/zap"
)

// Client represents a tweets client.
type Client struct {
	config Config
	logger *zap.Logger
}

// NewClient returns a new tweets client.
func NewClient(config Config, logger *zap.Logger) *Client {
	config.tweetsURL = fmt.Sprintf("http://%s/tweets/%s", config.Source, config.APIVersion)

	return &Client{
		config: config,
		logger: logger,
	}
}

// GetTweets returns tweets.
func (c *Client) GetTweets(ctx context.Context, user string, values v1.TweetsURLValues) (v1.TweetsResponse, error) {
	urlValues := url.Values{}
	err := schema.NewEncoder().Encode(&values, urlValues)
	if err != nil {
		return nil, err
	}

	userTweetsURL := fmt.Sprintf("%s/user/%s", c.config.tweetsURL, user)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, userTweetsURL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Form = urlValues

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.logger.Error("Closing response body failed", zap.Error(err))
		}
	}()

	var response v1.TweetsResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decoding tweets")
	}
	return response, nil
}
