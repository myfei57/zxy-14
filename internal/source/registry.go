package source

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"maskhub/internal/store"
)

var ErrNotFound = errors.New("source not found")

// Create registers a new data source for a tenant.
func Create(state *store.State, tenantID, name string) (*store.Source, error) {
	s := &store.Source{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      name,
		Status:    store.MaskPending,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutSource(s); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns a source by id.
func Get(state *store.State, id string) (*store.Source, error) {
	s, ok := state.Source(id)
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

// List returns all sources.
func List(state *store.State) []*store.Source {
	return state.Sources()
}
