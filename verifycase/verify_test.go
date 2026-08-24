package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"maskhub/internal/policy"
	"maskhub/internal/store"
)

func TestPolicyPublishWaitsForRuleSnapshotDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	p, err := policy.Create(state, "t1", "privacy")
	if err != nil {
		t.Fatal(err)
	}
	if err := state.PutRule(&store.Rule{ID: "r1", PolicyID: p.ID, Category: "phone", Pattern: "138", Replacer: "***", Order: 1}); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "snapshots"), 0o755); err != nil {
		t.Fatal(err)
	}
	blocker := filepath.Join(dir, "snapshots", p.ID+"-1.json")
	if err := os.RemoveAll(blocker); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blocker, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := policy.Publish(state, p.ID); err == nil {
		t.Fatal("expected publish to fail while the snapshot write is blocked")
	}
	after, _ := state.Policy(p.ID)
	if after.Status == store.PolicyActive {
		t.Fatal("policy became active although the rule snapshot was not durable")
	}
}
