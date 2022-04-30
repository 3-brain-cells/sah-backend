package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// Make sure Provider implements db.Provider
var _ db.Provider = &Provider{}

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

func (p *Provider) GetSingle(ctx context.Context, eventID string) (*types.Event, error) {
	collection := p.events()

	result := collection.FindOne(ctx, bson.M{"id": eventID})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, db.NewNotFoundError(eventID)
	}

	var event types.Event
	err := result.Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (p *Provider) CreatePartial(ctx context.Context, event types.Event) error {
	collection := p.events()
	// Ensure nested maps are empty
	if event.UserAvailability == nil {
		event.UserAvailability = make(map[string]types.UserAvailability)
	}
	if event.UserLocations == nil {
		event.UserLocations = make(map[string]types.UserLocation)
	}
	if event.UserVotes == nil {
		event.UserVotes = make(map[string]types.UserVotes)
	}
	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		if writeException, ok := err.(mongo.WriteException); ok && isDuplicate(writeException) {
			return db.NewDuplicateIDError(event.EventID)
		}
		return err
	}
	return nil
}

// Update updates an existing event
// Ignore the following fields:
// - creatorID
// - guildID
// - channelID
// - populated
// - voteOptions
// - userVotes
// - userAvailability
// - userLocations
// - userVotes
func (p *Provider) PopulateEvent(ctx context.Context, event types.Event, userID string) error {
	collection := p.events()

	// serialize to a string JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// deserialize that string into a map[string]interface{}
	mmap := make(map[string]interface{})
	err = json.Unmarshal(eventJSON, &mmap)
	if err != nil {
		return err
	}

	updateDocument := bson.D{}
	for k, v := range mmap {
		// Don't update the following fields:
		// - EventID
		// - CreatorID
		// - GuildID
		// - Populated
		// - UserVotes
		if k == "id" || k == "creator_id" || k == "guild_id" || k == "populated" || k == "user_votes" || k == "channel_id" || k == "user_availability" || k == "user_locations" || k == "vote_options" || k == "user_votes" {
			continue
		}
		updateDocument = append(updateDocument, bson.E{Key: k, Value: v})
	}

	filter := bson.D{{Key: "id", Value: event.EventID}}
	updateQuery := bson.D{{Key: "$set", Value: updateDocument}}

	_, err = collection.UpdateOne(ctx, filter, updateQuery)

	if err != nil {
		if writeException, ok := err.(mongo.WriteException); ok && isDuplicate(writeException) {
			return db.NewDuplicateIDError(event.EventID)
		}
		return err
	}
	return nil
}

// Update updates an existing event
func (p *Provider) PostVotes(ctx context.Context, userID string, votes types.UserVotes, eventID string) error {
	collection := p.events()

	votesJson, err := toRawRepresentation(votes)
	if err != nil {
		return fmt.Errorf("failed to marshal votes: %w", err)
	}

	filter := bson.D{{Key: "id", Value: eventID}}
	updateQuery := bson.M{
		"$set": bson.M{
			fmt.Sprintf("user_votes.%s", userID): rawToBson(votesJson),
		},
	}

	// print out using log. the userid and votes and event id and filter and update query
	log.Printf("userid: %#v", userID)
	log.Printf("votes: %#v", votesJson)
	log.Printf("event id: %#v", eventID)
	log.Printf("filter: %#v", filter)
	log.Printf("update query: %#v", updateQuery)

	_, err = collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update votes for userID=%s eventID=%s: %w", userID, eventID, err)
	}

	return nil
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

// Detects if the given write exception is caused by (in part)
// by a duplicate key error
func isDuplicate(writeException mongo.WriteException) bool {
	for _, writeError := range writeException.WriteErrors {
		if writeError.Code == duplicateError {
			return true
		}
	}

	return false
}

func toRawRepresentation(val interface{}) (map[string]interface{}, error) {
	jsonString, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value '%#v' to create raw representation: %w", val, err)
	}

	rawRepresentation := make(map[string]interface{})
	err = json.Unmarshal(jsonString, &rawRepresentation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw representation of JSON string '%s': %w", jsonString, err)
	}

	return rawRepresentation, nil
}

func rawToBson(rawRepresentation map[string]interface{}) bson.D {
	bsonRepresentation := bson.D{}
	for k, v := range rawRepresentation {
		bsonRepresentation = append(bsonRepresentation, bson.E{Key: k, Value: v})
	}
	return bsonRepresentation
}

func (p *Provider) PutUserAvailabilityAndLocation(ctx context.Context, userID string,
	availability types.UserAvailability, location types.UserLocation, eventID string) error {

	collection := p.events()

	availabilityJson, err := toRawRepresentation(availability)
	if err != nil {
		return fmt.Errorf("failed to marshal availability: %w", err)
	}

	locationJson, err := toRawRepresentation(location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	filter := bson.D{{Key: "id", Value: eventID}}
	updateQuery := bson.M{
		"$set": bson.M{
			fmt.Sprintf("user_availability.%s", userID): rawToBson(availabilityJson),
			fmt.Sprintf("user_locations.%s", userID):   rawToBson(locationJson),
		},
	}

	_, err = collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update availability and location for userID=%s eventID=%s: %w", userID, eventID, err)
	}

	return nil
}

func (p *Provider) GetAllEvents(ctx context.Context) ([]*types.Event, error) {
	collection := p.events()
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	var events []*types.Event
	for cursor.Next(ctx) {
		var event types.Event
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (p *Provider) UpdateVoteOptions(ctx context.Context, voteOptions types.VoteOption, eventID string) error {
	collection := p.events()

	voteOptionsJson, err := toRawRepresentation(voteOptions)
	if err != nil {
		return fmt.Errorf("failed to marshal voteOptions: %w", err)
	}

	filter := bson.D{{Key: "id", Value: eventID}}
	updateQuery := bson.M{
		"$set": bson.M{
			"vote_options": rawToBson(voteOptionsJson),
		},
	}

	_, err = collection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update vote options eventID=%s: %w", eventID, err)
	}

	return nil
}
