package classify

import (
	"errors"

	"maskhub/internal/store"
)

var ErrModel = errors.New("classification model not found")

// ApplyModel uses the current model snapshot to classify a field name.
func ApplyModel(state *store.State, name string) (string, error) {
	m, err := Current(state)
	if err != nil {
		return "", ErrModel
	}
	if m.Snapshot == "" {
		return "", ErrModel
	}
	if contains(name, "phone") || contains(name, "mobile") {
		return "phone", nil
	}
	if contains(name, "id") || contains(name, "card") {
		return "id_card", nil
	}
	if contains(name, "email") {
		return "email", nil
	}
	return "normal", nil
}

// SnapshotVersion returns the version of the current model.
func SnapshotVersion(state *store.State) (int, error) {
	m, err := Current(state)
	if err != nil {
		return 0, ErrModel
	}
	return m.Version, nil
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
