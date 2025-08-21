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
	// create the uri address
	uriAddress := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	// create the options
	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	// if use auth, set the auth
	if cfg.UseAuth {
		opt = opt.SetAuth(
			options.Credential{Username: cfg.User, Password: cfg.Password},
		)
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	// if enable tracing, add tracing
	if cfg.EnableTracing {
		opt.Monitor = otelmongo.NewMonitor()
	}

	// setup  https://github.com/Kamva/mgm
	err = mgm.SetDefaultConfig(nil, cfg.Database, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}
