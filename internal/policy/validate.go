package policy

import (
	"errors"

	"maskhub/internal/rule"
	"maskhub/internal/store"
)

var (
	ErrNoRules  = errors.New("policy has no rules")
	ErrNoSnapshot = errors.New("policy has no published snapshot")
)

// ValidateDraft checks that a draft policy has usable rules before publish.
func ValidateDraft(state *store.State, id string) error {
	p, ok := state.Policy(id)
	if !ok {
		return ErrNotFound
	}
	if p.Status != store.PolicyDraft && p.Status != store.PolicyPublished {
		return ErrState
	}
	rules := state.RulesForPolicy(id)
	if len(rules) == 0 {
		return ErrNoRules
	}
	for _, r := range rules {
		if r.Category == "" || r.Pattern == "" {
			return ErrNoRules
		}
	}
	return nil
}

// SnapshotForExecution returns the active policy's rule snapshot for masking.
func SnapshotForExecution(state *store.State, policyID string) ([]*store.Rule, error) {
	p, ok := state.Policy(policyID)
	if !ok {
		return nil, ErrNotFound
	}
	if p.Status != store.PolicyActive {
		return nil, ErrState
	}
	snap, ok := state.Snapshot(p.ID, p.Version)
	if !ok {
		return nil, ErrNoSnapshot
	}
	var rules []*store.Rule
	for _, id := range snap.RuleIDs {
		r, ok := state.Rule(id)
		if !ok {
			return nil, rule.ErrSnapshot
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// EffectiveRules resolves the current rules for a policy by version.
func EffectiveRules(state *store.State, policyID string, version int) ([]*store.Rule, error) {
	snap, ok := state.Snapshot(policyID, version)
	if !ok {
		return nil, ErrNoSnapshot
	}
	var rules []*store.Rule
	for _, id := range snap.RuleIDs {
		r, ok := state.Rule(id)
		if !ok {
			return nil, ErrNotFound
		}
		rules = append(rules, r)
	}
	return rules, nil
}
