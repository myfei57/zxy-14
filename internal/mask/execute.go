package mask

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"maskhub/internal/policy"
	"maskhub/internal/quota"
	"maskhub/internal/store"
)

var (
	ErrNotFound = errors.New("masking job not found")
	ErrState    = errors.New("masking job is not in the required state")
	ErrQuota    = errors.New("tenant quota exhausted")
)

// Execute starts a masking job for a source under a policy.
func Execute(state *store.State, tenantID, sourceID, policyID string, rows int) (*store.MaskJob, error) {
	p, err := policy.Get(state, policyID)
	if err != nil {
		return nil, err
	}
	size := int64(rows) * 1024
	if err := quota.Check(state, tenantID, size); err != nil {
		return nil, ErrQuota
	}
	j := &store.MaskJob{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		SourceID:  sourceID,
		PolicyID:  policyID,
		Status:    store.MaskRunning,
		Rows:      rows,
		QuotaUsed: size,
		CreatedAt: time.Now().UTC(),
	}
	if err := state.PutJob(j); err != nil {
		return nil, err
	}
	if err := quota.Reserve(state, tenantID, size); err != nil {
		return nil, err
	}
	_ = p
	return j, nil
}

// Finish marks a job done after its results are committed.
func Finish(state *store.State, id string) (*store.MaskJob, error) {
	j, ok := state.Job(id)
	if !ok {
		return nil, ErrNotFound
	}
	if j.Status != store.MaskRunning {
		return nil, ErrState
	}
	j.Status = store.MaskDone
	j.FinishedAt = time.Now().UTC()
	return j, state.PutJob(j)
}

// Fail marks a job failed.
func Fail(state *store.State, id string) error {
	j, ok := state.Job(id)
	if !ok {
		return ErrNotFound
	}
	j.Status = store.MaskFailed
	j.FinishedAt = time.Now().UTC()
	return state.PutJob(j)
}
