package mask

import (
	"time"

	"github.com/google/uuid"

	"maskhub/internal/policy"
	"maskhub/internal/store"
)

// Apply processes one row of a job and commits the masked result.
func Apply(state *store.State, jobID, rowID, category, value string) (*store.MaskResult, error) {
	j, ok := state.Job(jobID)
	if !ok {
		return nil, ErrNotFound
	}
	p, err := policy.Get(state, j.PolicyID)
	if err != nil {
		return nil, err
	}
	if p.Status != store.PolicyActive {
		return nil, policy.ErrState
	}
	if err := policy.Commit(state, p); err != nil {
		return nil, err
	}
	// Resolve rules from the published snapshot, not the live rule set, so a
	// job always runs against the verified snapshot captured at publish time.
	rules, err := policy.SnapshotForExecution(state, j.PolicyID)
	if err != nil {
		return nil, err
	}
	buf := NewBuffer()
	output := buf.Process(rules, category, value)
	res := &store.MaskResult{
		ID:        uuid.NewString(),
		JobID:     jobID,
		RowID:     rowID,
		Output:    output,
		Committed: true,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutResult(res); err != nil {
		return nil, err
	}
	return res, nil
}
