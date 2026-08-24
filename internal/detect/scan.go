package detect

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"maskhub/internal/classify"
	"maskhub/internal/store"
)

var ErrNotFound = errors.New("field not found")

// Scan detects a field and classifies it.
func Scan(state *store.State, sourceID, name string) (*store.Field, error) {
	f := &store.Field{
		ID:        uuid.NewString(),
		SourceID:  sourceID,
		Name:      name,
		Status:    store.FieldPending,
		ClassifiedAt: time.Now().UTC(),
	}
	if err := state.PutField(f); err != nil {
		return nil, err
	}
	if err := classify.ClassifyField(state, f); err != nil {
		return nil, err
	}
	return f, nil
}

// Fields returns all fields of a source.
func Fields(state *store.State, sourceID string) []*store.Field {
	return state.FieldsForSource(sourceID)
}

// Get returns a field by id.
func Get(state *store.State, id string) (*store.Field, error) {
	f, ok := state.Field(id)
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

// MarkClassified flips a field to classified after its result is durable.
func MarkClassified(state *store.State, id string) error {
	f, ok := state.Field(id)
	if !ok {
		return ErrNotFound
	}
	f.Status = store.FieldClassified
	f.ClassifiedAt = time.Now().UTC()
	return state.PutField(f)
}
