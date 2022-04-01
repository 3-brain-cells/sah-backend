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
	// GetAll returns all events
	GetAll(ctx context.Context, guildID string) ([]types.Event, error)

	// GetSingle returns a single event
	GetSingle(ctx context.Context, eventID int) (types.Event, error)

	// Create creates a new event
	Create(ctx context.Context, event types.Event) error

	// Update updates an existing event
	Update(ctx context.Context, event types.Event) error

	// Delete deletes an existing event
	Delete(ctx context.Context, eventID int) error
}
