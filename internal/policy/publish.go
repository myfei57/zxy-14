package policy

import (
	"time"

	"maskhub/internal/rule"
	"maskhub/internal/store"
)

// Publish activates a draft policy. The rule snapshot is durably written and
// verified to resolve to a complete rule set BEFORE the active switch, so a
// failed or partial snapshot leaves the policy in its pre-publish state and
// masking never runs against a partial rule set.
func Publish(state *store.State, id string) (*store.Policy, error) {
	p, ok := state.Policy(id)
	if !ok {
		return nil, ErrNotFound
	}
	if p.Status != store.PolicyDraft && p.Status != store.PolicyPublished {
		return nil, ErrState
	}
	// Persist the rule snapshot first. Only once the snapshot has landed on
	// disk and been verified complete do we flip the policy active.
	snap, err := rule.SnapshotRules(state, p.ID, p.Version)
	if err != nil {
		return nil, err
	}
	if _, err := EffectiveRules(state, snap.PolicyID, snap.Version); err != nil {
		return nil, err
	}
	p.Status = store.PolicyActive
	p.PublishedAt = time.Now().UTC()
	p.SnapshotID = snap.PolicyID
	if err := state.PutPolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}
