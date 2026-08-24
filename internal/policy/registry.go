package policy

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"maskhub/internal/store"
)

var (
	ErrNotFound = errors.New("policy not found")
	ErrState    = errors.New("policy is not in the required state")
)

// Create starts a draft policy.
func Create(state *store.State, tenantID, name string) (*store.Policy, error) {
	p := &store.Policy{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      name,
		Version:   1,
		Status:    store.PolicyDraft,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutPolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Get returns a policy by id.
func Get(state *store.State, id string) (*store.Policy, error) {
	p, ok := state.Policy(id)
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// List returns all policies.
func List(state *store.State) []*store.Policy {
	return state.Policies()
}

// Active returns the currently active policy of a tenant.
func Active(state *store.State, tenantID string) (*store.Policy, error) {
	for _, p := range state.Policies() {
		if p.TenantID == tenantID && p.Status == store.PolicyActive {
			return p, nil
		}
	}
	return nil, ErrNotFound
}
