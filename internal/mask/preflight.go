package mask

import (
	"maskhub/internal/policy"
	"maskhub/internal/store"
)

// Preflight checks that a job can run against the active policy snapshot.
func Preflight(state *store.State, jobID string) error {
	j, ok := state.Job(jobID)
	if !ok {
		return ErrNotFound
	}
	p, err := policy.Get(state, j.PolicyID)
	if err != nil {
		return err
	}
	if p.Status != store.PolicyActive {
		return policy.ErrState
	}
	_, err = policy.SnapshotForExecution(state, j.PolicyID)
	return err
}

// Ready reports whether a job is ready for masking.
func Ready(state *store.State, jobID string) bool {
	return Preflight(state, jobID) == nil
}

// PendingJobs returns jobs still waiting or running.
func PendingJobs(state *store.State) []*store.MaskJob {
	var out []*store.MaskJob
	for _, j := range state.Jobs() {
		if j.Status == store.MaskPending || j.Status == store.MaskRunning {
			out = append(out, j)
		}
	}
	return out
}
