package policy

import (
	"maskhub/internal/store"
)

// Commit durably stores a policy edit and returns the committed policy.
func Commit(state *store.State, p *store.Policy) error {
	if p == nil {
		return ErrNotFound
	}
	return state.PutPolicy(p)
}

// Snapshot returns the durable snapshot of the active policy.
func Snapshot(state *store.State, policyID string, version int) (string, error) {
	snap, err := loadSnapshot(state, policyID, version)
	if err != nil {
		return "", err
	}
	return snap.PolicyID, nil
}

func loadSnapshot(state *store.State, policyID string, version int) (*store.Snapshot, error) {
	s, ok := state.Snapshot(policyID, version)
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}
