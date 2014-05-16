package slog

import (
	"sort"
	"testing"
)

// We steal the "level" field to store the expected sorted order
var testingRules = rules{
	{"hello", 4},
	{"world", 0},
	{"omg", 2},
	{"helloworld", 3},
	{"hell", 5},
	{"omgponies", 1},
	{"", 6},
}

func TestRuleSorting(t *testing.T) {
	t.Parallel()
	sort.Sort(testingRules)
	for i, rule := range testingRules {
		if i != int(rule.level) {
			t.Errorf("Rule %q was in pos. %d, expected %d",
				rule.selector, i, rule.level)
		}
	}
}
