package verifycase

import (
	"testing"

	"maskhub/internal/mask"
	"maskhub/internal/policy"
	"maskhub/internal/store"
)

func TestMaskAppliesAfterPolicyCommit(t *testing.T) {
	state := store.NewState(t.TempDir())
	p, err := policy.Create(state, "t1", "privacy")
	if err != nil {
		t.Fatal(err)
	}
	if err := state.PutRule(&store.Rule{ID: "r1", PolicyID: p.ID, Category: "phone", Pattern: "138", Replacer: "***", Order: 1}); err != nil {
		t.Fatal(err)
	}
	// The policy edit is never committed (still draft), so masking must refuse.
	j := &store.MaskJob{ID: "j1", TenantID: "t1", PolicyID: p.ID, Status: store.MaskRunning}
	if err := state.PutJob(j); err != nil {
		t.Fatal(err)
	}
	if _, err := mask.Apply(state, "j1", "row1", "phone", "1381234"); err == nil {
		t.Fatal("masking applied rules from an uncommitted policy edit")
	}
}
