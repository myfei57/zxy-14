package policy

import "maskhub/internal/store"

// Rules returns the rules of a policy ordered by their order field.
func Rules(state *store.State, policyID string) []*store.Rule {
	return state.RulesForPolicy(policyID)
}

// RuleCount returns the number of rules attached to a policy.
func RuleCount(state *store.State, policyID string) int {
	return len(state.RulesForPolicy(policyID))
}

// CategoryCoverage returns which field categories a policy covers.
func CategoryCoverage(state *store.State, policyID string) []string {
	seen := map[string]bool{}
	var out []string
	for _, r := range state.RulesForPolicy(policyID) {
		if !seen[r.Category] {
			seen[r.Category] = true
			out = append(out, r.Category)
		}
	}
	return out
}
