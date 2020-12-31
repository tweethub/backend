package mongodb

import (
	"context"

	"github.com/tweethub/backend/cmd/tweets/storage"
	"github.com/tweethub/backend/pkg/database"
	"github.com/tweethub/backend/pkg/database/mongodb"
	"github.com/tweethub/backend/pkg/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// MongoDB wraps database access.
type MongoDB struct {
	conn   *mongo.Client
	db     *mongo.Database
	logger *zap.Logger
}

// Init returns a new mongodb twitter storage.
func Init(ctx context.Context, cfg database.Config, logger *zap.Logger) storage.Twitter {
	logger.Info(service.InitDatabase)

	conn, err := mongodb.Connect(ctx, cfg)
	if err != nil {
		logger.Fatal(service.InitDatabaseFailed, zap.Error(err))
	}

	return &MongoDB{
		conn:   conn,
		db:     conn.Database(cfg.Name),
		logger: logger,
	}
}

// Close closes the database connection.
func (mng *MongoDB) Close(ctx context.Context) {
	if err := mng.conn.Disconnect(ctx); err != nil {
		mng.logger.Error("Closing database connection failed",
			zap.Error(err))
	}
}
