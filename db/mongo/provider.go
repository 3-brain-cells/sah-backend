package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/env"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/rs/zerolog"
)

const (
	duplicateError = 11000
)

// Provider implements the Provider interface for a MongoDB connection
type Provider struct {
	logger        zerolog.Logger
	connectionURI string
	databaseName  string
	clusterName   string
	client        *mongo.Client
}

// NewProvider creates a new provider and loads values in from the environment
func NewProvider(logger zerolog.Logger) (*Provider, error) {
	username, err := env.GetEnv("MongoDB username", "MONGO_DB_USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := env.GetEnv("MongoDB password", "MONGO_DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	clusterName, err := env.GetEnv("MongoDB cluster name ", "MONGO_DB_CLUSTER_NAME")
	if err != nil {
		return nil, err
	}

	databaseName, err := env.GetEnv("MongoDB database name ", "MONGO_DB_DATABASE_NAME")
	if err != nil {
		return nil, err
	}

	connectionURI := fmt.Sprintf("mongodb+srv://%s:%s@%s.6ta2w.mongodb.net/%s?retryWrites=true&w=majority",
		username, password, clusterName, databaseName)
	return &Provider{
		logger:        logger,
		connectionURI: connectionURI,
		databaseName:  databaseName,
		clusterName:   clusterName,
		client:        nil,
	}, nil
}

// Connect connects to the MongoDB server and creates any indices as necessary
func (p *Provider) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(p.connectionURI))
	if err != nil {
		return err
	}

	// Ping the primary
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return err
	}

	p.client = client

	// Initialize any collections/indices
	err = p.initialize(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Disconnect terminates the connection to the MongoDB server
func (p *Provider) Disconnect(ctx context.Context) error {
	err := p.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Create anything needed for the database,
// like indices
func (p *Provider) initialize(ctx context.Context) error {
	p.logger.
		Info().
		Str("database_name", p.databaseName).
		Str("cluster_name", p.clusterName).
		Msg("initializing the MongoDB database")

	_, err := p.events().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) events() *mongo.Collection {
	return p.client.Database(p.databaseName).Collection("events")
}

// GetEvent given its ID
func (p *Provider) GetEvent(ctx context.Context, id string) (*types.Event, error) {
	collection := p.events()
	result := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, db.NewNotFoundError(id)
	}

	var event types.Event
	err := result.Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
