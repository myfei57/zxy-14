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
	res.Committed = true
	if err := state.PutResult(res); err != nil {
		return err
	}
	return audit.Record(state, "mask_result", res.ID, res.Output)
}
