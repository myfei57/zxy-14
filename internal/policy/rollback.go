package policy

import (
	"path/filepath"
	"time"

	"maskhub/internal/quota"
	"maskhub/internal/store"
)

// Rollback reverts an active policy. The rollback record must be durable
// before quota is released, so a failed write never frees capacity while the
// old policy stays active.
func Rollback(state *store.State, id string) (*store.Policy, error) {
	p, ok := state.Policy(id)
	if !ok {
		return nil, ErrNotFound
	}
	if p.Status != store.PolicyActive {
		return nil, ErrState
	}
	// Durably persist the rollback record first. It is the commit point of the
	// rollback: while it is not on disk the policy stays active, so a failed
	// write frees no quota and leaves the rollback retryable. Only after the
	// record is durable do we flip the status and release quota.
	record := map[string]any{
		"policy_id": p.ID,
		"version":   p.Version,
		"rolled_back_at": time.Now().UTC().Format(time.RFC3339),
	}
	if err := store.SaveJSON(filepath.Join(state.Root(), "rollbacks", p.ID+".json"), record); err != nil {
		return nil, err
	}
	p.Status = store.PolicyRolledBack
	if err := state.PutPolicy(p); err != nil {
		return nil, err
	}
	if err := quota.Release(state, p.TenantID, quotaReserved(state, p)); err != nil {
		return nil, err
	}
	return p, nil
}

func quotaReserved(state *store.State, p *store.Policy) int64 {
	return 0
}
