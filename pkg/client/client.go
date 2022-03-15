package client

import (
	"context"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

var ProviderClientSet = wire.NewSet(NewClient)

func NewClient(ctx context.Context, log *zap.Logger) *mongo.Client {
	uri := viper.GetString("MONGODB_URI")

	log.Info("Connecting to MongoDB", zap.String("uri", uri))

	clientOpts := options.Client().
		ApplyURI(uri)

	// Set connect timeout to 15 seconds
	ctxConn, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create new client and connect to MongoDB
	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		log.Panic("Failed to connect to mongodb", zap.Error(err))
	}

	// Set the ping timeout to 5 seconds
	ctxPing, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	// Ping the primary
	if err = client.Ping(ctxPing, readpref.Primary()); err != nil {
		log.Panic("Failed to Ping() to mongodb", zap.Error(err))
	}

	log.Info("Successfully connected and pinged.")

	return client
}
