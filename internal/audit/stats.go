package audit

import "maskhub/internal/store"

// ByAction groups audit entries by action name.
func ByAction(state *store.State) map[string]int {
	out := make(map[string]int)
	for _, e := range state.AuditEntries() {
		out[e.Action]++
	}
	return out
}

// Total returns the total number of audit entries.
func Total(state *store.State) int {
	return len(state.AuditEntries())
}
