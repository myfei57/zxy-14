package mask

import (
	"maskhub/internal/rule"
	"maskhub/internal/store"
)

// Buffer is the per-row redaction buffer. It must be reset between rows so an
// unmatched row never inherits the previous row's replacement token.
type Buffer struct {
	value string
	token string
	used  bool
}

// NewBuffer creates an empty per-row buffer.
func NewBuffer() *Buffer {
	return &Buffer{}
}

// Reset clears the buffer for a new row.
func (b *Buffer) Reset() {
	b.value = ""
	b.token = ""
	b.used = false
}

// Process applies rules to a row value, using the buffer for replacements.
func (b *Buffer) Process(rules []*store.Rule, category, value string) string {
	b.Reset()
	b.value = value
	r := rule.Match(rules, category)
	if r == nil {
		return value
	}
	out := rule.Apply(r, value)
	if out != value {
		b.token = r.Replacer
		b.used = true
	}
	return out
}

// Token returns the replacement token stored in the buffer.
func (b *Buffer) Token() string {
	return b.token
}

// Used reports whether the buffer holds a replacement token.
func (b *Buffer) Used() bool {
	return b.used
}
