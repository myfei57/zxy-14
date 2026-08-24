package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"maskhub/internal/policy"
	"maskhub/internal/store"
)

func TestRollbackReleasesQuotaAfterDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	p, err := policy.Create(state, "t1", "privacy")
	if err != nil {
		t.Fatal(err)
	}
	if err := state.PutRule(&store.Rule{ID: "r1", PolicyID: p.ID, Category: "phone", Pattern: "138", Replacer: "***", Order: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := policy.Publish(state, p.ID); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "rollbacks"), 0o755); err != nil {
		t.Fatal(err)
	}
	blocker := filepath.Join(dir, "rollbacks", p.ID+".json")
	if err := os.RemoveAll(blocker); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blocker, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := policy.Rollback(state, p.ID); err == nil {
		t.Fatal("expected rollback to fail while the rollback record write is blocked")
	}
	after, _ := state.Policy(p.ID)
	if after.Status == store.PolicyRolledBack {
		t.Fatal("policy rolled back although the rollback record was not durable")
	}
}
