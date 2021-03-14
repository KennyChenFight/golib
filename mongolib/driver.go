// Package mongolib is for encapsulating go.mongodb.org/mongo-driver any operations
//
// As a quick start:
// 	ctx := context.Background()
//	mgo, err := mongolib.NewMongoDriver(ctx, mongolib.Config{
//		AppName: "test",
//		URL:         "mongodb://localhost:27017",
//		Database:    "golang",
//	}, nil)
//	if err != nil {
//		panic(err)
//	}
//	collection := mgo.Database.Collection("example")
//  collection.find(...)
package mongolib

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	AppName     string
	URL         string
	Database    string
	MaxPoolSize uint64
}

func NewMongoDriver(ctx context.Context, cfg Config, opts ...*options.ClientOptions) (*MongoDriver, error) {
	option := options.Client()
	option.SetAppName(cfg.AppName)
	option.ApplyURI(cfg.URL)

	if cfg.MaxPoolSize != 0 {
		option.SetMaxPoolSize(cfg.MaxPoolSize)
	}

	var (
		client *mongo.Client
		err    error
	)
	if opts != nil {
		opts = append(opts, option)
		client, err = mongo.NewClient(opts...)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = mongo.NewClient(option)
		if err != nil {
			return nil, err
		}
	}

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDriver{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

type MongoDriver struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (m *MongoDriver) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
