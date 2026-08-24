package verifycase

import (
	"testing"

	"maskhub/internal/detect"
	"maskhub/internal/store"
)

func TestDetectOriginRefreshesAfterReclassify(t *testing.T) {
	state := store.NewState(t.TempDir())
	f := &store.Field{ID: "f1", SourceID: "s1", Name: "mobile", Category: "normal", Status: store.FieldClassified}
	if err := state.PutField(f); err != nil {
		t.Fatal(err)
	}
	// A reclassification changes the category; detection must refresh its origin.
	f.Category = "phone"
	if err := state.PutField(f); err != nil {
		t.Fatal(err)
	}
	if err := detect.RefreshOrigin(state, "f1"); err != nil {
		t.Fatal(err)
	}
	after, _ := state.Field("f1")
	if after.Category != "phone" {
		t.Fatalf("category %q, want phone after refresh", after.Category)
	}
	if after.Status != store.FieldReclassified {
		t.Fatalf("status %q, want reclassified", after.Status)
	}
}
