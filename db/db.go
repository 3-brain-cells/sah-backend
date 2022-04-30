package db

import (
	"context"

	"github.com/3-brain-cells/sah-backend/types"
)

// Provider represents a database provider implementation
type Provider interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error

	EventProvider
}

// EventProvider provides CRUD operations for types.Event structs
type EventProvider interface {
	// GetSingle returns a single event
	GetSingle(ctx context.Context, eventID string) (*types.Event, error)

	// Create creates a new partial event (before it is populated)
	CreatePartial(ctx context.Context, event types.Event) error

	// pass in a partial event struct
	// ignore the following fields:
	// - creatorID
	// - guildID
	// - channelID
	// - populated
	// - voteOptions
	// - userVotes
	// If userID is not the creator ID of the event, an error is returned.
	PopulateEvent(ctx context.Context, event types.Event, userID string) error

	// Update updates an existing event
	PostVotes(ctx context.Context, userID string, votes types.UserVotes, eventID string) error

	// PutUserAvailabilityAndLocation updates the user availability and location
	PutUserAvailabilityAndLocation(ctx context.Context, userID string, availability types.UserAvailability, location types.UserLocation, eventID string) error

	// GetAllEvents returns all events in the database
	GetAllEvents(ctx context.Context) ([]*types.Event, error)

	// UpdateVoteOptions updates the vote options for times and locations
	UpdateVoteOptions(ctx context.Context, voteOptions types.VoteOption, eventID string) error

	// Delete deletes an existing event
	// Delete(ctx context.Context, eventID int) error
}
