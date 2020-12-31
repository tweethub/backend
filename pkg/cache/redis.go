package cache

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// ConnectRedis connects to redis.
func ConnectRedis(conf *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "connecting to redis")
	}
	return client, nil
}
