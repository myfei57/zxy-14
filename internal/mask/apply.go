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
	// Apply only the rules from the active policy's durable snapshot, so
	// masking never runs against uncommitted edits, save failures, or draft
	// rules. Mirrors mask.Batch, which already uses this path.
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
