package quota

import (
	"errors"

	"maskhub/internal/store"
)

var (
	ErrNotFound = errors.New("tenant not found")
	ErrExceeded = errors.New("tenant quota exceeded")
)

// Check verifies that adding size stays within the tenant quota.
func Check(state *store.State, tenantID string, size int64) error {
	if size < 0 {
		return nil
	}
	used := state.QuotaUsed(tenantID)
	limit := state.QuotaLimit(tenantID)
	if limit > 0 && used+size > limit {
		return ErrExceeded
	}
	return nil
}

// Used returns the current used quota of a tenant.
func Used(state *store.State, tenantID string) int64 {
	return state.QuotaUsed(tenantID)
}

// Reserve adds size to the tenant quota.
func Reserve(state *store.State, tenantID string, size int64) error {
	return state.AddQuota(tenantID, size)
}

// Release returns size from the tenant quota.
func Release(state *store.State, tenantID string, size int64) error {
	return state.AddQuota(tenantID, -size)
}
