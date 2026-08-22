package policy

import (
	"time"

	"maskhub/internal/rule"
	"maskhub/internal/store"
)

// Publish activates a draft policy. The rule snapshot must be durable before
// the active switch, so masking never runs against a partial rule set.
func Publish(state *store.State, id string) (*store.Policy, error) {
	p, ok := state.Policy(id)
	if !ok {
		return nil, ErrNotFound
	}
	if p.Status != store.PolicyDraft && p.Status != store.PolicyPublished {
		return nil, ErrState
	}
	snap, err := rule.SnapshotRules(state, p.ID, p.Version)
	if err != nil {
		return nil, err
	}
	p.SnapshotID = snap.PolicyID
	p.Status = store.PolicyActive
	p.PublishedAt = time.Now().UTC()
	if err := state.PutPolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}
