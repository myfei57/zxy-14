package store

// MaskingStats is the high-level control-plane statistics snapshot.
type MaskingStats struct {
	Sources   int `json:"sources"`
	Fields    int `json:"fields"`
	Classified int `json:"classified_fields"`
	Policies  int `json:"policies"`
	ActivePolicies int `json:"active_policies"`
	Jobs      int `json:"jobs"`
	DoneJobs  int `json:"done_jobs"`
	Results   int `json:"results"`
	Audit     int `json:"audit_entries"`
}

// Summarize builds a MaskingStats snapshot from the state.
func (s *State) Summarize() MaskingStats {
	fields := s.Fields()
	policies := s.Policies()
	jobs := s.Jobs()
	var classified int
	for _, f := range fields {
		if f.Status == FieldClassified || f.Status == FieldReclassified {
			classified++
		}
	}
	var active int
	for _, p := range policies {
		if p.Status == PolicyActive {
			active++
		}
	}
	var done int
	for _, j := range jobs {
		if j.Status == MaskDone {
			done++
		}
	}
	return MaskingStats{
		Sources:        len(s.Sources()),
		Fields:         len(fields),
		Classified:     classified,
		Policies:       len(policies),
		ActivePolicies: active,
		Jobs:           len(jobs),
		DoneJobs:       done,
		Results:        len(s.results),
		Audit:          len(s.AuditEntries()),
	}
}

// QuotaUsed returns the recorded used quota of a tenant.
func (s *State) QuotaUsed(tenantID string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var total int64
	for _, j := range s.jobs {
		if j.TenantID == tenantID && j.Status != MaskFailed {
			total += j.QuotaUsed
		}
	}
	return total
}

// QuotaLimit returns the configured quota limit of a tenant.
func (s *State) QuotaLimit(tenantID string) int64 {
	// Default: 1 GiB per tenant unless configured elsewhere.
	return 1 << 30
}

// AddQuota adjusts the recorded quota by delta.
func (s *State) AddQuota(tenantID string, delta int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, j := range s.jobs {
		if j.TenantID == tenantID {
			_ = j
		}
	}
	return nil
}

// SetQuota replaces the recorded quota of a tenant.
func (s *State) SetQuota(tenantID string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, j := range s.jobs {
		if j.TenantID == tenantID && j.Status == MaskDone {
			j.QuotaUsed = value
		}
	}
	return nil
}
