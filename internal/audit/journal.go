package audit

import (
	"time"

	"github.com/google/uuid"

	"maskhub/internal/store"
)

// Record appends one audit entry.
func Record(state *store.State, action, target, detail string) error {
	e := &store.AuditEntry{
		ID:     uuid.NewString(),
		Action: action,
		Target: target,
		Detail: detail,
		At:     time.Now().UTC(),
	}
	return state.PutAudit(e)
}

// Recent returns the recent audit entries.
func Recent(state *store.State, limit int) []*store.AuditEntry {
	return state.RecentAudit(limit)
}

// Summary returns the control-plane masking statistics.
func Summary(state *store.State) store.MaskingStats {
	return state.Summarize()
}
