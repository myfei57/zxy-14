package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"maskhub/internal/mask"
	"maskhub/internal/store"
)

func TestAuditAfterMaskCommitDurable(t *testing.T) {
	dir := t.TempDir()
	state := store.NewState(dir)
	j := &store.MaskJob{ID: "j1", TenantID: "t1", PolicyID: "p1", Status: store.MaskRunning}
	if err := state.PutJob(j); err != nil {
		t.Fatal(err)
	}
	res := &store.MaskResult{ID: "r1", JobID: "j1", RowID: "row1", Output: "***", Committed: false}
	if err := state.PutResult(res); err != nil {
		t.Fatal(err)
	}
	// Block the result file so committing the output must fail.
	resultFile := filepath.Join(dir, "results", "r1.json")
	if err := os.RemoveAll(resultFile); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(resultFile, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := mask.Commit(state, res); err == nil {
		t.Fatal("expected commit to fail while the result write is blocked")
	}
	// No audit entry may exist for output that was never durably committed.
	if len(state.AuditEntries()) != 0 {
		t.Fatal("audit recorded masking output that was never committed")
	}
}
