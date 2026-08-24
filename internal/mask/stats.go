package mask

import "maskhub/internal/store"

// JobStats summarizes one masking job.
type JobStats struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Rows      int    `json:"rows"`
	Results   int    `json:"results"`
	Committed int    `json:"committed"`
	QuotaUsed int64  `json:"quota_used"`
}

// StatsOf computes the statistics of a masking job.
func StatsOf(state *store.State, jobID string) (JobStats, error) {
	j, ok := state.Job(jobID)
	if !ok {
		return JobStats{}, ErrNotFound
	}
	results := state.ResultsForJob(jobID)
	s := JobStats{
		JobID:     jobID,
		Status:    j.Status,
		Rows:      j.Rows,
		Results:   len(results),
		QuotaUsed: j.QuotaUsed,
	}
	for _, r := range results {
		if r.Committed {
			s.Committed++
		}
	}
	return s, nil
}

// ListStats returns statistics for all jobs.
func ListStats(state *store.State) []JobStats {
	var out []JobStats
	for _, j := range state.Jobs() {
		s, err := StatsOf(state, j.ID)
		if err == nil {
			out = append(out, s)
		}
	}
	return out
}
