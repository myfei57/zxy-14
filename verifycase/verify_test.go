package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"maskhub/internal/classify"
	"maskhub/internal/source"
	"maskhub/internal/store"
)

func TestDetectWaitsForClassificationDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	src, err := source.Create(state, "t1", "orders")
	if err != nil {
		t.Fatal(err)
	}
	f := &store.Field{ID: "f1", SourceID: src.ID, Name: "mobile", Status: store.FieldPending}
	if err := state.PutField(f); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "classifications"), 0o755); err != nil {
		t.Fatal(err)
	}
	blocker := filepath.Join(dir, "classifications", "f1.json")
	if err := os.RemoveAll(blocker); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(blocker, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := classify.ClassifyField(state, f); err == nil {
		t.Fatal("expected classify to fail while the classification write is blocked")
	}
	after, _ := state.Field("f1")
	if after.Status == store.FieldClassified {
		t.Fatal("field marked classified although the classification result was not durable")
	}
}
