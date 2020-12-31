package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tweethub/backend/pkg/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect returns a connected mongodb instance.
func Connect(ctx context.Context, cfg database.Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating new database")
	}

	if err = client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "connecting to database server")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "pinging database server")
	}
	return client, nil
}
