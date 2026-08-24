package verifycase

import (
	"testing"

	"maskhub/internal/classify"
	"maskhub/internal/store"
)

func TestDetectionUsesCurrentModelAfterUpdate(t *testing.T) {
	state := store.NewState(t.TempDir())
	if _, err := classify.Update(state, "v1"); err != nil {
		t.Fatal(err)
	}
	if _, err := classify.Update(state, "v2"); err != nil {
		t.Fatal(err)
	}
	version, err := classify.SnapshotVersion(state)
	if err != nil {
		t.Fatal(err)
	}
	if version != 2 {
		t.Fatalf("model version %d, want 2 after update", version)
	}
	category, err := classify.ApplyModel(state, "mobile_number")
	if err != nil {
		t.Fatal(err)
	}
	if category != "phone" {
		t.Fatalf("category %q, want phone", category)
	}
}
