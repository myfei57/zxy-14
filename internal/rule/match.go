package rule

import (
	"strconv"
	"strings"
	"time"

	"maskhub/internal/store"
)

// Match returns the rule that applies to a category, or nil.
func Match(rules []*store.Rule, category string) *store.Rule {
	for _, r := range rules {
		if r.Category == category {
			return r
		}
	}
	return nil
}

// Apply replaces matching content in a value using the rule.
func Apply(rule *store.Rule, value string) string {
	if rule == nil {
		return value
	}
	if rule.Pattern == "" || !strings.Contains(value, rule.Pattern) {
		return value
	}
	return strings.ReplaceAll(value, rule.Pattern, rule.Replacer)
}

func itoa(v int) string {
	return formatInt(v)
}

func formatInt(v int) string {
	return strconv.Itoa(v)
}

func timeNow() time.Time {
	return time.Now().UTC()
}
