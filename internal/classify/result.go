package classify

import (
	"errors"
	"path/filepath"
	"time"

	"maskhub/internal/store"
)

var errNoModel = errors.New("no classification model registered")

// ClassifyField durably stores the classification result, then marks the field.
func ClassifyField(state *store.State, f *store.Field) error {
	if f == nil {
		return ErrFieldRequired
	}
	if err := detectMark(state, f.ID); err != nil {
		return err
	}
	result := map[string]any{
		"field_id": f.ID,
		"name":     f.Name,
		"category": f.Category,
		"at":       time.Now().UTC().Format(time.RFC3339),
	}
	// Durably record the classification before flipping the field state.
	if err := store.SaveJSON(filepath.Join(state.Root(), "classifications", f.ID+".json"), result); err != nil {
		return err
	}
	return nil
}

// Reclassify updates a field's category and stores the new result durably.
func Reclassify(state *store.State, f *store.Field) error {
	if f == nil {
		return ErrFieldRequired
	}
	result := map[string]any{
		"field_id": f.ID,
		"name":     f.Name,
		"category": f.Category,
		"reclassified": true,
		"at":       time.Now().UTC().Format(time.RFC3339),
	}
	if err := store.SaveJSON(filepath.Join(state.Root(), "classifications", f.ID+".json"), result); err != nil {
		return err
	}
	f.Status = store.FieldReclassified
	f.ClassifiedAt = time.Now().UTC()
	return state.PutField(f)
}

var ErrFieldRequired = errors.New("field is required")

func detectMark(state *store.State, fieldID string) error {
	f, ok := state.Field(fieldID)
	if !ok {
		return ErrFieldRequired
	}
	f.Status = store.FieldClassified
	f.ClassifiedAt = time.Now().UTC()
	return state.PutField(f)
}
