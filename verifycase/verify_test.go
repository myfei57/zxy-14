package verifycase

import (
	"testing"

	"maskhub/internal/mask"
	"maskhub/internal/policy"
	"maskhub/internal/quota"
	"maskhub/internal/store"
)

func TestMaskingChecksQuotaBeforeTransform(t *testing.T) {
	state := store.NewState(t.TempDir())
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
	src := "s1"
	if err := state.PutSource(&store.Source{ID: src, TenantID: "t1", Name: "orders"}); err != nil {
		t.Fatal(err)
	}
	// Exhaust the quota before starting.
	if err := state.PutJob(&store.MaskJob{ID: "j0", TenantID: "t1", Status: store.MaskDone, QuotaUsed: 1 << 40}); err != nil {
		t.Fatal(err)
	}
	if _, err := mask.Execute(state, "t1", src, p.ID, 10); err == nil {
		t.Fatal("expected masking to fail when the tenant quota is exhausted")
	}
	_ = quota.Used
}
