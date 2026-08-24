package rule

import (
	"errors"
	"path/filepath"

	"maskhub/internal/store"
)

var ErrSnapshot = errors.New("rule snapshot write failed")

// SnapshotRules durably stores the rule snapshot of a policy version.
func SnapshotRules(state *store.State, policyID string, version int) (*store.Snapshot, error) {
	rules := state.RulesForPolicy(policyID)
	ids := make([]string, 0, len(rules))
	for _, r := range rules {
		ids = append(ids, r.ID)
	}
	snap := &store.Snapshot{
		PolicyID: policyID,
		Version:  version,
		RuleIDs:  ids,
		CreatedAt: timeNow(),
	}
	// Durable write must succeed before the snapshot is usable.
	if err := store.SaveJSON(filepath.Join(state.Root(), "snapshots", policyID+"-"+itoa(version)+".json"), snap); err != nil {
		return nil, ErrSnapshot
	}
	if err := state.PutSnapshot(snap); err != nil {
		return nil, err
	}
	return snap, nil
}

// LoadSnapshot returns the snapshot of a policy version.
func LoadSnapshot(state *store.State, policyID string, version int) (*store.Snapshot, error) {
	snap, ok := state.Snapshot(policyID, version)
	if !ok {
		return nil, ErrSnapshot
	}
	return snap, nil
}
