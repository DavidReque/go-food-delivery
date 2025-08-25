package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	// connectTimeout is the timeout for the connection to the database
	connectTimeout = 60 * time.Second // 60 seconds
	// maxConnIdleTime is the maximum idle time for a connection
	maxConnIdleTime = 3 * time.Minute // 3 minutes
	// minPoolSize is the minimum number of connections in the pool
	minPoolSize = 20 // 20 connections
	// maxPoolSize is the maximum number of connections in the pool
	maxPoolSize = 300 // 300 connections
)

// NewMongoDB Create new MongoDB client
func NewMongoDB(cfg *MongoDbOptions) (*mongo.Client, error) {

	// Apply default values if empty (temporary configuration problem)
	if cfg.Host == "" && !cfg.UseAtlas {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 27017
	}
	if cfg.Database == "" {
		cfg.Database = "catalogs_read_service"
	}

	var uriAddress string

	// First check if Atlas should be used
	if cfg.UseAtlas {
		if cfg.AtlasURI == "" {
			return nil, fmt.Errorf("MongoDB Atlas URI is empty but useAtlas is true")
		}
		uriAddress = cfg.AtlasURI
	} else {
		// Use local configuration
		if cfg.UseAuth {
			uriAddress = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
		} else {
			uriAddress = fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
		}
	}

	opt := options.Client().
		ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	if cfg.EnableTracing {
		opt.Monitor = otelmongo.NewMonitor()
	}

	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	// setup  https://github.com/Kamva/mgm
	// This will fail if cfg.Database is empty, providing a more accurate error source.
	if err = mgm.SetDefaultConfig(nil, cfg.Database, opt); err != nil {
		return nil, fmt.Errorf("failed to set mgm default config: %w", err)
	}

	return client, nil
}
