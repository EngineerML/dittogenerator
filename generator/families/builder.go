// Package families holds the case families our generator composes.
//
// A family is a deterministic function that seeds a few memory records and
// returns one challenge case grounded in them. The Builder owns the haystack
// (waves, pair and session IDs, timestamps, renderers, noise) so families
// only describe content: what is stated, in which wave, and what is asked.
//
// Every case a family emits must pass the pack contract: its gold answer
// appears in the named evidence pairs, not in the question, and the haystack
// holds the decoys listed as distractors, so an agent that dumps all of
// memory fails.
//
// All randomness flows from the request seed through Builder.Rng; nothing
// reads the clock or the environment, so a request always yields the same
// bytes (the scorer re-runs a seed and rejects a generator that drifts).
package families

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// Waves is the number of seeding waves. A chain spreads its steps over
// distinct waves, so a change and the note that points back at it sit
// sessions apart.
const Waves = 8

// Path is the reading path a case's question and records steer an agent to.
// The top agents route on surface features: a long identifier carried by
// several records sends a question to their case-file (program) path, a
// phrase shared with several multi-paragraph notes to their topic path, and
// a full name or nothing at all to their records lane. Each path reads
// chains with different rules, so the same chain is a different test on
// each.
type Path int

const (
	PathRecords Path = iota
	PathProgram
	PathTopic
)

func (p Path) String() string {
	return [...]string{"records", "program", "topic"}[p]
}

// Builder accumulates a pack's haystack and cases.
type Builder struct {
	Rng    *rand.Rand
	waves  []protocol.SeedRequest
	cases  []challenge.Case
	used   map[string]bool
	tokens map[string]bool
	epoch  time.Time
}

// NewBuilder returns an empty builder for one request.
func NewBuilder(req challenge.GenerateRequest) *Builder {
	b := &Builder{
		Rng:    req.Rand(),
		waves:  make([]protocol.SeedRequest, Waves),
		used:   map[string]bool{},
		tokens: map[string]bool{},
		epoch:  time.Date(2026, time.March, 2, 8, 0, 0, 0, time.UTC),
	}
	for i := range b.waves {
		b.waves[i].Wave = i
	}
	return b
}

// token returns a fresh lowercase alphanumeric string of length n, unique
// within the pack.
func (b *Builder) token(n int) string {
	const alphabet = "abcdefghijkmnpqrstuvwxyz23456789"
	for {
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = alphabet[b.Rng.Intn(len(alphabet))]
		}
		if s := string(buf); !b.tokens[s] {
			b.tokens[s] = true
			return s
		}
	}
}

// Say seeds one conversation pair in a wave and returns its pair ID. IDs and
// sessions are opaque, so neither gives away which records belong together.
func (b *Builder) Say(wave int, prompt, response string) string {
	id := "m" + b.token(9)
	b.waves[wave].Pairs = append(b.waves[wave].Pairs, protocol.MemoryPair{
		PairID:    id,
		SessionID: fmt.Sprintf("s%d%s", wave, b.token(4)),
		Prompt:    prompt,
		Response:  response,
	})
	return id
}

// Add appends a case. Cases run after the last wave so every fact is seeded.
func (b *Builder) Add(c challenge.Case) {
	c.RunAfterWave = Waves - 1
	b.cases = append(b.cases, c)
}

// Pick draws n distinct values from pool that no earlier family has used, so
// two cases never collide on a value (which would make one case's gold
// another's distractor).
func (b *Builder) Pick(pool []string, n int) []string {
	var out []string
	for _, i := range b.Rng.Perm(len(pool)) {
		if len(out) == n {
			break
		}
		if v := pool[i]; !b.used[v] {
			b.used[v] = true
			out = append(out, v)
		}
	}
	if len(out) < n {
		panic(fmt.Sprintf("value pool exhausted: need %d more distinct values", n-len(out)))
	}
	return out
}

// Any draws one value from pool without reserving it. Use it only for parts
// of a question that a reserved value already makes unique.
func (b *Builder) Any(pool []string) string {
	return pool[b.Rng.Intn(len(pool))]
}

// Person returns a reserved full name. First names and surnames are drawn
// separately, so the pack never runs short of people, and both halves are
// plain capitalised words the agents' name patterns read.
func (b *Builder) Person() string {
	for {
		name := b.Any(firstNames) + " " + b.Any(surnames)
		if !b.used[name] {
			b.used[name] = true
			return name
		}
	}
}

// People returns n reserved full names.
func (b *Builder) People(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = b.Person()
	}
	return out
}

// CaseFile returns a fresh case-file identifier: twelve characters with at
// least one letter and one digit, the shape the agents treat as a record key.
func (b *Builder) CaseFile() string {
	for {
		id := b.token(12)
		letter, digit := false, false
		for _, r := range id {
			if r >= '0' && r <= '9' {
				digit = true
			} else {
				letter = true
			}
		}
		if letter && digit {
			return id
		}
	}
}

// Spread returns k distinct waves in increasing order, the waves a chain of
// k steps is told in.
func (b *Builder) Spread(k int) []int {
	w := b.Rng.Perm(Waves)[:k]
	sort.Ints(w)
	return w
}

// Filler seeds an unrelated pair in a random wave, so a case's facts are not
// next to each other in the haystack.
func (b *Builder) Filler() {
	line := filler[b.Rng.Intn(len(filler))]
	b.Say(b.Rng.Intn(Waves), line[0], line[1])
}

// Len is the number of cases added so far.
func (b *Builder) Len() int { return len(b.cases) }

// Pack assembles the challenge pack bound to the request.
func (b *Builder) Pack(req challenge.GenerateRequest, name, version string) challenge.Pack {
	// Timestamps are laid out once every pair is known: wave by wave, 20 to
	// 80 minutes apart, with a day between waves. Families interleave waves as
	// they seed facts, so stamping at Say time could put an earlier wave after
	// a later one.
	at := b.epoch
	for w := range b.waves {
		for p := range b.waves[w].Pairs {
			at = at.Add(time.Duration(20+b.Rng.Intn(61)) * time.Minute)
			b.waves[w].Pairs[p].Timestamp = at.Format(time.RFC3339)
		}
		at = at.Add(24 * time.Hour)
	}
	return challenge.Pack{
		Contract:            challenge.ContractV1,
		GradingBenchVersion: req.GradingBenchVersion,
		MinerHotkey:         req.MinerHotkey,
		Generator:           challenge.Generator{Name: name, Version: version, Seed: req.Seed},
		MemoryWaves:         b.waves,
		Cases:               b.cases,
	}
}
