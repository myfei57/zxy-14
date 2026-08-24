package schedule

import (
	"time"

	"maskhub/internal/quota"
	"maskhub/internal/store"
)

// Runner drives the periodic maintenance loop of the masking gateway.
type Runner struct {
	state   *store.State
	stop    chan struct{}
	done    chan struct{}
	options Options
}

// Options controls the maintenance cadence.
type Options struct {
	ReconcileEvery time.Duration
}

// DefaultOptions returns the default maintenance cadence.
func DefaultOptions() Options {
	return Options{ReconcileEvery: 15 * time.Second}
}

// NewRunner creates a maintenance runner for the state.
func NewRunner(state *store.State, opts Options) *Runner {
	return &Runner{
		state:   state,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		options: opts,
	}
}

// Start launches the maintenance loop in a background goroutine.
func (r *Runner) Start() {
	go r.loop()
}

// Stop signals the loop to exit and waits for it.
func (r *Runner) Stop() {
	close(r.stop)
	<-r.done
}

func (r *Runner) loop() {
	defer close(r.done)
	ticker := time.NewTicker(r.options.ReconcileEvery)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			r.runReconcile()
		}
	}
}

func (r *Runner) runReconcile() {
	tenants := map[string]bool{}
	for _, j := range r.state.Jobs() {
		tenants[j.TenantID] = true
	}
	for tenantID := range tenants {
		_ = quota.Reconcile(r.state, tenantID)
	}
}
