package detect

import (
	"maskhub/internal/store"
)

// Origin returns the current classification origin for a field.
func Origin(state *store.State, fieldID string) (*store.Field, error) {
	f, ok := state.Field(fieldID)
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

// RefreshOrigin reloads a field's classification after a reclassification.
func RefreshOrigin(state *store.State, fieldID string) error {
	_, ok := state.Field(fieldID)
	if !ok {
		return ErrNotFound
	}
	return nil
}
