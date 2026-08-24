package mask

import (
	"maskhub/internal/audit"
	"maskhub/internal/store"
)

// Commit durably stores a masked result and then records the audit entry.
func Commit(state *store.State, res *store.MaskResult) error {
	if res == nil {
		return ErrNotFound
	}
	if err := audit.Record(state, "mask_result", res.ID, res.Output); err != nil {
		return err
	}
	res.Committed = true
	return state.PutResult(res)
}
