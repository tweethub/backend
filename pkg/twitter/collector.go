package twitter

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/tweethub/backend/pkg/json"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Message struct {
	Err    error
	Tweets ArchiveResponse
}

type MessageStream chan Message

// Config represents collector configuration.
type Config struct {
	// CollectPastYearsTweets will make the collector to collect past year tweets too.
	CollectPastYearsTweets bool
	CollectingInterval     time.Duration
}

func (cfg Config) Validate() error {
	var err error

	if cfg.CollectingInterval <= 0 {
		err = multierr.Append(err,
			errors.New("collecting interval should be bigger than 0 in collector configuration"))
	}

	return err
}

// Collector represents a collector of tweets.
type Collector struct {
	logger *zap.Logger
	config Config
}

// NewCollector returns new collector.
func NewCollector(config Config, logger *zap.Logger) *Collector {
	return &Collector{
		logger: logger,
		config: config,
	}
}

// CollectOverTime collects tweets over time.
func (col *Collector) CollectOverTime(ctx context.Context) MessageStream {
	col.logger.Info("Collecting tweets",
		zap.String("time-interval", col.config.CollectingInterval.String()),
		zap.Bool("collect-past-year-tweets", col.config.CollectPastYearsTweets))

	msgStream := make(MessageStream)

	go func() {
		ticker := time.NewTicker(col.config.CollectingInterval)
		for {
			select {
			case <-ctx.Done():
				close(msgStream)
				return
			case <-ticker.C:
				tweets, err := col.CollectTweets()
				msgStream <- Message{
					Tweets: tweets,
					Err:    err,
				}
			}
		}
	}()
	return msgStream
}

// CollectTweets collects and returns tweets.
func (col *Collector) CollectTweets() (ArchiveResponse, error) {
	var allTweets ArchiveResponse

	urls := col.archiveURLs()
	for _, url := range urls {
		resp, err := col.getArchive(url)
		if err != nil {
			return nil, err
		}

		allTweets = append(allTweets, resp...)
	}
	return allTweets, nil
}

// getArchive collects tweets from a source archive.
func (col *Collector) getArchive(source string) (ArchiveResponse, error) {
	resp, err := http.Get(source) // nolint:gosec,noctx
	if err != nil {
		return nil, errors.Wrap(err, "http get request")
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			col.logger.Error("Closing response body failed", zap.Error(err))
		}
	}()

	var response ArchiveResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "decoding tweets")
	}
	return response, err
}

// archiveURLs returns the URLs for the tweets archives.
func (col *Collector) archiveURLs() []string {
	archiveYear := UserFirstTweetDate().Year()
	currentYear := time.Now().Year()
	urls := make([]string, 0, (currentYear-archiveYear)+1)

	for ; archiveYear <= currentYear; archiveYear++ {
		urls = append(urls, fmt.Sprintf(sourceURLScheme, archiveYear))
	}

	if !col.config.CollectPastYearsTweets {
		sort.Strings(urls)
		urls = urls[len(urls)-1:] // Get only the current year archive URL.
	}
	return urls
}
