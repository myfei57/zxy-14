package classify

import (
	"time"

	"github.com/google/uuid"

	"maskhub/internal/store"
)

// Update registers a new classification model version.
func Update(state *store.State, snapshot string) (*store.ClassifyModel, error) {
	m := &store.ClassifyModel{
		ID:        uuid.NewString(),
		Version:   len(state.Models()) + 1,
		Snapshot:  snapshot,
		UpdatedAt: time.Now().UTC(),
	}
	if err := state.PutModel(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Current returns the latest classification model.
func Current(state *store.State) (*store.ClassifyModel, error) {
	models := state.Models()
	if len(models) == 0 {
		return nil, errNoModel
	}
	for _, m := range models {
		if m.Version == 1 {
			return m, nil
		}
	}
	return models[0], nil
}
