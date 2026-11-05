package families

import (
	"regexp"
	"strings"
	"testing"
)

// maxRounds is the number of rounds of the thirteen default slots the
// generator runs at its MaxCases (80): seven, the last one partial.
const maxRounds = 7

// TestPoolsFillMaxCases holds the reserved pools to what MaxCases draws.
// Builder.Pick panics on an exhausted pool, so an undersized pool would fail
// every request near the top of the range.
func TestPoolsFillMaxCases(t *testing.T) {
	need := map[string]struct {
		pool    []string
		perPack int
	}{
		"cities":     {cities, 22},
		"projects":   {projects, 9},
		"milestones": {milestones, 5},
		"statuses":   {statuses, 3},
		"firms":      {firms, 2},
		"domains":    {domains, 2},
		"nicknames":  {nicknames, 5},
		"shopNames":  {shopNames, 8},
	}
	for name, n := range need {
		if got, want := len(n.pool), n.perPack*maxRounds; got < want {
			t.Errorf("%s has %d values; %d rounds need %d", name, got, maxRounds, want)
		}
	}
}

// TestPoolValuesStayApart refuses a value that appears, as whole words, inside
// another value of any pool. Such a pair lets one case's gold answer ground or
// decoy another case's question.
func TestPoolValuesStayApart(t *testing.T) {
	pools := [][]string{cities, firstNames, surnames, nicknames, projects, firms, statuses, milestones, events, threads, rosters, schedules, shopNames}
	var all []string
	seen := map[string]bool{}
	for _, pool := range pools {
		for _, v := range pool {
			key := strings.ToLower(v)
			if seen[key] {
				t.Errorf("value %q appears twice", v)
			}
			seen[key] = true
			all = append(all, v)
		}
	}
	for _, a := range all {
		word := regexp.MustCompile(`(?i)(^|[^a-z])` + regexp.QuoteMeta(a) + `($|[^a-z])`)
		for _, b := range all {
			if a != b && word.MatchString(b) {
				t.Errorf("%q contains %q", b, a)
			}
		}
	}
}
