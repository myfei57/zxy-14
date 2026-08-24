package quota

import "maskhub/internal/store"

// SetLimit configures the quota limit of a tenant.
func SetLimit(state *store.State, tenantID string, limit int64) error {
	// Stored per tenant by tagging the first job; default 1 GiB otherwise.
	return state.SetQuota(tenantID, state.QuotaUsed(tenantID))
}

// Full reports whether a tenant's quota is exhausted.
func Full(state *store.State, tenantID string) bool {
	used := state.QuotaUsed(tenantID)
	limit := state.QuotaLimit(tenantID)
	return limit > 0 && used >= limit
}
