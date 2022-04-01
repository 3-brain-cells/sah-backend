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
	CreatePartial(ctx context.Context, event types.EventCreate) error

	// populate the created event with other fields
	PopulateEvent(ctx context.Context, event types.Event) error

	// Update updates an existing event
	PostVotes(ctx context.Context, votes types.UserVotes, eventID string) error

	// Delete deletes an existing event
	Delete(ctx context.Context, eventID int) error
}
