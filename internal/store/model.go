package store

import "time"

// State machine constants shared across the masking control plane.
const (
	FieldPending     = "pending"
	FieldClassified  = "classified"
	FieldReclassified = "reclassified"

	PolicyDraft     = "draft"
	PolicyPublished = "published"
	PolicyActive    = "active"
	PolicyRolledBack = "rolled_back"

	MaskPending = "pending"
	MaskRunning = "masking"
	MaskDone    = "done"
	MaskFailed  = "failed"
)

// Source is one data source feeding the masking gateway.
type Source struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Field is one detected field of a source.
type Field struct {
	ID           string    `json:"id"`
	SourceID     string    `json:"source_id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Status       string    `json:"status"`
	OriginVersion int      `json:"origin_version"`
	ClassifiedAt time.Time `json:"classified_at"`
}

// Policy is one masking policy version.
type Policy struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	Status      string    `json:"status"`
	SnapshotID  string    `json:"snapshot_id"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
}

// Rule is one masking rule bound to a field category.
type Rule struct {
	ID        string    `json:"id"`
	PolicyID  string    `json:"policy_id"`
	Category  string    `json:"category"`
	Pattern   string    `json:"pattern"`
	Replacer  string    `json:"replacer"`
	Order     int       `json:"order"`
}

// Snapshot is the durable rule snapshot of a published policy.
type Snapshot struct {
	PolicyID  string   `json:"policy_id"`
	Version   int      `json:"version"`
	RuleIDs   []string `json:"rule_ids"`
	CreatedAt time.Time `json:"created_at"`
}

// MaskJob is one masking execution task.
type MaskJob struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	SourceID  string    `json:"source_id"`
	PolicyID  string    `json:"policy_id"`
	Status    string    `json:"status"`
	Rows      int       `json:"rows"`
	QuotaUsed int64     `json:"quota_used"`
	CreatedAt time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}

// MaskResult is one row's masked output.
type MaskResult struct {
	ID        string    `json:"id"`
	JobID     string    `json:"job_id"`
	RowID     string    `json:"row_id"`
	Output    string    `json:"output"`
	Committed bool      `json:"committed"`
	CreatedAt time.Time `json:"created_at"`
}

// ClassifyModel is one classification model version.
type ClassifyModel struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	Snapshot  string    `json:"snapshot"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditEntry records one control-plane operation.
type AuditEntry struct {
	ID     string    `json:"id"`
	Action string    `json:"action"`
	Target string    `json:"target"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}
