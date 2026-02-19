package store

import (
	"context"

	"github.com/devaloi/htmxapp/internal/model"
)

// ContactStore defines the interface for contact persistence.
type ContactStore interface {
	List(ctx context.Context, search string) ([]model.Contact, error)
	Get(ctx context.Context, id string) (model.Contact, error)
	Create(ctx context.Context, c model.Contact) (model.Contact, error)
	Update(ctx context.Context, c model.Contact) (model.Contact, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) int
}
