package mask

import (
	"maskhub/internal/policy"
	"maskhub/internal/store"
)

// RowInput is one row to be masked.
type RowInput struct {
	RowID    string `json:"row_id"`
	Category string `json:"category"`
	Value    string `json:"value"`
}

// Batch processes multiple rows for a job and returns the committed results.
func Batch(state *store.State, jobID string, rows []RowInput) ([]*store.MaskResult, error) {
	j, ok := state.Job(jobID)
	if !ok {
		return nil, ErrNotFound
	}
	rules, err := policy.SnapshotForExecution(state, j.PolicyID)
	if err != nil {
		return nil, err
	}
	buf := NewBuffer()
	var out []*store.MaskResult
	for _, row := range rows {
		output := buf.Process(rules, row.Category, row.Value)
		res := &store.MaskResult{
			ID:        newResultID(),
			JobID:     jobID,
			RowID:     row.RowID,
			Output:    output,
			Committed: true,
			CreatedAt: nowUTC(),
		}
		if err := state.PutResult(res); err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, nil
}

// Results returns the results of a job.
func Results(state *store.State, jobID string) []*store.MaskResult {
	return state.ResultsForJob(jobID)
}

func newResultID() string {
	return randomID()
}

func nowUTC() timeType {
	return timeNow()
}
