package quota

import "maskhub/internal/store"

// Ledger returns the recorded used bytes of a tenant.
func Ledger(state *store.State, tenantID string) int64 {
	return state.QuotaUsed(tenantID)
}

// Reconcile recomputes the tenant quota from masking job usage.
func Reconcile(state *store.State, tenantID string) error {
	var total int64
	for _, j := range state.Jobs() {
		if j.TenantID == tenantID && j.Status == store.MaskDone {
			total += j.QuotaUsed
		}
	}
	return state.SetQuota(tenantID, total)
}

// Remaining returns how many bytes a tenant may still process.
func Remaining(state *store.State, tenantID string) int64 {
	used := state.QuotaUsed(tenantID)
	limit := state.QuotaLimit(tenantID)
	if limit <= 0 {
		return 0
	}
	r := limit - used
	if r < 0 {
		return 0
	}
	return r
}
