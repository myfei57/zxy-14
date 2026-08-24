package store

import (
	"path/filepath"
	"sort"
	"sync"
)

// State is the central in-memory plus file-backed registry.
type State struct {
	mu        sync.RWMutex
	root      string
	sources   map[string]*Source
	fields    map[string]*Field
	policies  map[string]*Policy
	rules     map[string]*Rule
	snapshots map[string]*Snapshot
	jobs      map[string]*MaskJob
	results   map[string]*MaskResult
	models    map[string]*ClassifyModel
	audit     map[string]*AuditEntry
}

// NewState creates a State rooted at dir.
func NewState(dir string) *State {
	return &State{
		root:      dir,
		sources:   make(map[string]*Source),
		fields:    make(map[string]*Field),
		policies:  make(map[string]*Policy),
		rules:     make(map[string]*Rule),
		snapshots: make(map[string]*Snapshot),
		jobs:      make(map[string]*MaskJob),
		results:   make(map[string]*MaskResult),
		models:    make(map[string]*ClassifyModel),
		audit:     make(map[string]*AuditEntry),
	}
}

// Root returns the state root directory.
func (s *State) Root() string {
	return s.root
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *State) PutSource(v *Source) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sources[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "sources", v.ID+".json"), v)
}

func (s *State) Source(id string) (*Source, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sources[id]
	return v, ok
}

func (s *State) Sources() []*Source {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Source, 0, len(s.sources))
	for _, k := range sortedKeys(s.sources) {
		out = append(out, s.sources[k])
	}
	return out
}

func (s *State) PutField(v *Field) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fields[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "fields", v.ID+".json"), v)
}

func (s *State) Field(id string) (*Field, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.fields[id]
	return v, ok
}

func (s *State) Fields() []*Field {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Field, 0, len(s.fields))
	for _, k := range sortedKeys(s.fields) {
		out = append(out, s.fields[k])
	}
	return out
}

func (s *State) FieldsForSource(sourceID string) []*Field {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Field
	for _, f := range s.fields {
		if f.SourceID == sourceID {
			out = append(out, f)
		}
	}
	return out
}

func (s *State) PutPolicy(v *Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "policies", v.ID+".json"), v)
}

func (s *State) Policy(id string) (*Policy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.policies[id]
	return v, ok
}

func (s *State) Policies() []*Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Policy, 0, len(s.policies))
	for _, k := range sortedKeys(s.policies) {
		out = append(out, s.policies[k])
	}
	return out
}

func (s *State) PutRule(v *Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "rules", v.ID+".json"), v)
}

func (s *State) Rule(id string) (*Rule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.rules[id]
	return v, ok
}

func (s *State) RulesForPolicy(policyID string) []*Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Rule
	for _, r := range s.rules {
		if r.PolicyID == policyID {
			out = append(out, r)
		}
	}
	return out
}

func (s *State) PutSnapshot(v *Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots[v.PolicyID+"#"+itoa(v.Version)] = v
	return nil
}

func (s *State) Snapshot(policyID string, version int) (*Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.snapshots[policyID+"#"+itoa(version)]
	return v, ok
}

func (s *State) PutJob(v *MaskJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "jobs", v.ID+".json"), v)
}

func (s *State) Job(id string) (*MaskJob, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.jobs[id]
	return v, ok
}

func (s *State) Jobs() []*MaskJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*MaskJob, 0, len(s.jobs))
	for _, k := range sortedKeys(s.jobs) {
		out = append(out, s.jobs[k])
	}
	return out
}

func (s *State) PutResult(v *MaskResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "results", v.ID+".json"), v)
}

func (s *State) Result(id string) (*MaskResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.results[id]
	return v, ok
}

func (s *State) ResultsForJob(jobID string) []*MaskResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*MaskResult
	for _, r := range s.results {
		if r.JobID == jobID {
			out = append(out, r)
		}
	}
	return out
}

func (s *State) PutModel(v *ClassifyModel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.models[v.ID] = v
	return SaveJSON(filepath.Join(s.root, "models", v.ID+".json"), v)
}

func (s *State) Models() []*ClassifyModel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*ClassifyModel, 0, len(s.models))
	for _, k := range sortedKeys(s.models) {
		out = append(out, s.models[k])
	}
	return out
}

func (s *State) PutAudit(e *AuditEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audit[e.ID] = e
	return SaveJSON(filepath.Join(s.root, "audit", e.ID+".json"), e)
}

func (s *State) AuditEntries() []*AuditEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*AuditEntry, 0, len(s.audit))
	for _, k := range sortedKeys(s.audit) {
		out = append(out, s.audit[k])
	}
	return out
}

func (s *State) RecentAudit(limit int) []*AuditEntry {
	all := s.AuditEntries()
	if len(all) <= limit {
		return all
	}
	return all[len(all)-limit:]
}
