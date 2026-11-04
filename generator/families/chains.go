package families

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// Chain families: one value changes several times and a later note decides
// which one stands, often by pointing back without naming it. The agent
// research measured these as the band where gpt-oss-20b agents are right 20
// to 80 percent of the time: "reverts to the previous entry instead" was right
// 9 of 86 times, a membership fold asked "after the change" 10 of 18, a struck
// entry whose field is worded three ways split about 30 percent of groups.

var cityLines = Grammar{
	"set": {
		"{s} is going to be in {v}.",
		"We booked {s} in {v}.",
		"{s} will take place in {v}.",
		"Planning note: {s} is set for {v}.",
	},
	"change": {
		"Update on {s}: it moved to {v}.",
		"{s} is now in {v} instead.",
		"Change of plan for {s}, it's {v} now.",
		"{s} has been moved to {v}.",
		"Heads up, {s} shifted to {v}.",
	},
	"ask": {
		"Which city is {a} in now?",
		"Where does {a} take place now? Just the city.",
		"What city is {a} set for at this point?",
	},
}

// cityChain states cities[0] and then each later city in its own wave, with
// filler in between, and returns the pair IDs in order.
func cityChain(b *Builder, s Subject, cities []string, waves []int) []string {
	ids := []string{b.Record(s, waves[0], b.Phrase(cityLines, "set", "s", s.Name, "v", cities[0]), cities[0])}
	for k, city := range cities[1:] {
		b.Filler()
		ids = append(ids, b.Record(s, waves[k+1], b.Phrase(cityLines, "change", "s", s.Name, "v", city), city))
	}
	return ids
}

// valueCase is a single-value case. Its ID carries the family, variant and
// path, so a panel run can be read per family and path.
func valueCase(kind string, path Path, i int, question, gold string, distractors []string, evidence ...string) challenge.Case {
	return challenge.Case{
		MemoryCase: protocol.MemoryCase{
			ID:                fmt.Sprintf("%s-%s-%d", kind, path, i),
			QuestionType:      kind,
			Question:          question,
			ExpectedAnswer:    gold,
			DistractorAnswers: distractors,
		},
		EvidencePairIDs: evidence,
	}
}

// pathFor rotates a family's cases over the three paths, offset per family so
// one round puts different families on different paths.
func pathFor(i, offset int) Path { return Path((i + offset) % 3) }

var revertLines = Grammar{
	// The benchmark's own unnamed revert ends "; the newest entry governs",
	// and with it no local agent picked the previous city (0 of 18, panel
	// 2026-11-04): the clause reads as an order to take the newest one. The
	// mid forms keep the unnamed pointer without it; the soft ones are the
	// forms agents got right 87 to 100 percent of the time.
	"mid": {
		"From now on {s} reverts to the previous entry instead.",
		"{s} reverts to the previous entry.",
		"For {s}, go back to the previous entry.",
	},
	"soft": {
		"{s} goes back to the entry before the last change.",
		"Scrap the latest change to {s} and keep what we had right before it.",
		"{s} is back to what it was before the most recent change.",
		"We undid the last change to {s}.",
	},
	"original": {
		"{s} goes back to the city we had at the very start.",
		"Back to the original plan for {s}: the first city we picked.",
		"Forget all the moves, {s} returns to where it was first booked.",
		"{s} reverts to the first entry we ever made for it.",
	},
}

// RevertPrevious moves an event three or four times, then sends it back to
// the previous entry without naming it. The gold is the city before the last
// move.
func RevertPrevious(b *Builder, i int) challenge.Case {
	// Always the records lane: there the unnamed revert splits the agents (mid
	// 5 of 12, soft 10 of 12, panel 2026-11-04), while long topic notes bury
	// the revert line (mid 0 of 12) and case files make the soft form trivial.
	path := PathRecords
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 3+i%2)
	w := b.Spread(len(c) + 1)
	ids := cityChain(b, s, c, w)
	variant := "soft"
	if i%2 == 0 {
		variant = "mid"
	}
	ev := b.Record(s, w[len(w)-1], b.Phrase(revertLines, variant, "s", s.Name))
	gold := c[len(c)-2]
	return valueCase("revert-previous-"+variant, path, i, b.Phrase(cityLines, "ask", "a", s.Ask), gold, without(c, gold), ids[len(ids)-2], ids[len(ids)-1], ev)
}

// RevertOriginal moves an event two or three times, then sends it back to
// the first city without naming it.
func RevertOriginal(b *Builder, i int) challenge.Case {
	path := pathFor(i, 1)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 3+i%2)
	w := b.Spread(len(c) + 1)
	ids := cityChain(b, s, c, w)
	ev := b.Record(s, w[len(w)-1], b.Phrase(revertLines, "original", "s", s.Name))
	return valueCase("revert-original", path, i, b.Phrase(cityLines, "ask", "a", s.Ask), c[0], c[1:], ids[0], ev)
}

var struckLines = Grammar{
	// The same slot in three wordings: item identity across rewordings is
	// what the agents lose here.
	"w0": {"{s} is going to be in {v}.", "We booked {s} in {v}."},
	"w1": {"The venue for {s} is now in {v}.", "{s} will be hosted in {v} after all."},
	"w2": {"Location for {s}: {v}.", "{s} gets run out of {v} now."},
	"strike": {
		"The last entry about where {s} is held was logged in error; that entry no longer stands.",
		"Ignore the most recent location update for {s}, it was entered by mistake.",
		"Correction: the latest note on where {s} happens was a mistake, strike it.",
	},
}

// StruckEntry changes an event's city once or twice, each time wording the
// location differently, then strikes the latest entry. The gold is the city
// that held before it.
func StruckEntry(b *Builder, i int) challenge.Case {
	path := pathFor(i, 2)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 2+i%2)
	w := b.Spread(len(c) + 1)
	var ids []string
	for k, city := range c {
		ids = append(ids, b.Record(s, w[k], b.Phrase(struckLines, fmt.Sprintf("w%d", k), "s", s.Name, "v", city), city))
		b.Filler()
	}
	ev := b.Record(s, w[len(w)-1], b.Phrase(struckLines, "strike", "s", s.Name))
	gold := c[len(c)-2]
	return valueCase("struck-entry", path, i, b.Phrase(cityLines, "ask", "a", s.Ask), gold, without(c, gold), ids[len(ids)-2], ev)
}

var cueLines = Grammar{
	"plain":     {"{s} happens in {v}.", "{s} is in {v}."},
	"plainmove": {"{s} happens in {v} from here on.", "We're doing {s} in {v} going forward."},
	"draft": {
		"Draft update for sign-off: {s} moves to {v}.",
		"Update (pending approval): {s} is now {v}.",
	},
	"dropped": {
		"The draft about moving {s} was turned down, so nothing changed there.",
		"Sign-off came back negative on the {s} move; it stays put.",
	},
	"proposal": {
		"Update for review: {s} could move to {v}.",
		"Proposal on the table: {s} in {v}.",
	},
	"declined": {
		"That proposal for {s} was declined.",
		"We passed on the suggested move for {s}.",
	},
}

// MisplacedCue writes the governing change plainly and the change that never
// happened in update words: a draft "update" that sign-off turned down, or a
// plain move followed by a declined proposal.
func MisplacedCue(b *Builder, i int) challenge.Case {
	path := pathFor(i, 0)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 3)
	if i%2 == 0 {
		w := b.Spread(3)
		set := b.Record(s, w[0], b.Phrase(cueLines, "plain", "s", s.Name, "v", c[0]), c[0])
		b.Filler()
		b.Record(s, w[1], b.Phrase(cueLines, "draft", "s", s.Name, "v", c[1]), c[1])
		ev := b.Record(s, w[2], b.Phrase(cueLines, "dropped", "s", s.Name))
		// c[2] is the decoy another event of the same subject never had; seed
		// it as someone else's suggestion so memory holds it.
		b.Record(s, b.Rng.Intn(Waves), fmt.Sprintf("Someone floated %s for a future one, nothing more.", c[2]), c[2])
		return valueCase("misplaced-cue-draft", path, i, b.Phrase(cityLines, "ask", "a", s.Ask), c[0], c[1:], set, ev)
	}
	w := b.Spread(4)
	b.Record(s, w[0], b.Phrase(cueLines, "plain", "s", s.Name, "v", c[0]), c[0])
	move := b.Record(s, w[1], b.Phrase(cueLines, "plainmove", "s", s.Name, "v", c[1]), c[1])
	b.Filler()
	b.Record(s, w[2], b.Phrase(cueLines, "proposal", "s", s.Name, "v", c[2]), c[2])
	ev := b.Record(s, w[3], b.Phrase(cueLines, "declined", "s", s.Name))
	return valueCase("misplaced-cue-proposal", path, i, b.Phrase(cityLines, "ask", "a", s.Ask), c[1], []string{c[0], c[2]}, move, ev)
}

var memberLines = Grammar{
	"start":  {"{s} for now: {v}.", "Current names on {s}: {v}."},
	"remove": {"Take {v} off {s}.", "{v} dropped out of {s}."},
	"add":    {"Add {v} to {s}.", "{v} joins {s}."},
	"swap":   {"Swap {v} out for {w} on {s}.", "On {s}, {w} replaces {v}."},
	"readd":  {"Put {v} back on {s} after all.", "{v} is back on {s}."},
	"ask":    {"Who is on {a} after the change?"},
	"askall": {"Who is on {a} now, with every change applied?"},
}

// MembershipFold starts a list of four people and removes, adds and swaps
// members across waves; the hard variant puts a removed member back. The
// question either presupposes a single change ("after the change") or asks
// for every change applied.
func MembershipFold(b *Builder, i int) challenge.Case {
	path := pathFor(i, 1)
	s := b.Subject(path, b.Any(rosters))
	m := b.People(6)
	readd := i%2 == 0
	steps := 4
	if readd {
		steps = 5
	}
	w := b.Spread(steps)
	ids := []string{
		b.Record(s, w[0], b.Phrase(memberLines, "start", "s", s.Name, "v", list(m[:4])), m[:4]...),
		b.Record(s, w[1], b.Phrase(memberLines, "remove", "s", s.Name, "v", m[1]), m[1]),
		b.Record(s, w[2], b.Phrase(memberLines, "add", "s", s.Name, "v", m[4]), m[4]),
		b.Record(s, w[3], b.Phrase(memberLines, "swap", "s", s.Name, "v", m[2], "w", m[5]), m[2], m[5]),
	}
	final, wrong := []string{m[0], m[3], m[4], m[5]}, []string{m[1], m[2]}
	if readd {
		ids = append(ids, b.Record(s, w[4], b.Phrase(memberLines, "readd", "s", s.Name, "v", m[1]), m[1]))
		final, wrong = []string{m[0], m[1], m[3], m[4], m[5]}, []string{m[2]}
	}
	ask := "askall"
	if i%3 != 2 {
		ask = "ask"
	}
	kind := "membership-fold"
	if readd {
		kind += "-readd"
	}
	return challenge.Case{
		MemoryCase: protocol.MemoryCase{
			ID:                fmt.Sprintf("%s-%s-%d", kind, path, i),
			QuestionType:      "membership-fold",
			Question:          b.Phrase(memberLines, ask, "a", s.Ask),
			ExpectedAnswer:    strings.Join(final, ", "),
			AnswerKind:        protocol.AnswerList,
			AnswerItems:       final,
			DistractorAnswers: wrong,
		},
		EvidencePairIDs: ids,
	}
}

var orderLines = Grammar{
	"abs": {"For {s}, the {e} is on {d}.", "{s}: {e} booked for {d}."},
	"rel": {"For {s}, the {e} is {k} days {dir} the {r}.", "{s}: the {e} lands {k} days {dir} the {r}."},
	"ask": {
		"Which milestone on {a} comes right before the {t}?",
		"On {a}, which milestone comes just before the {t}?",
	},
}

// EventOrder dates five milestones of one schedule, two of them only relative
// to another ("three days after the tasting"), files them in random waves,
// and asks for the one right before a later milestone. The answer is never
// the first or last milestone.
func EventOrder(b *Builder, i int) challenge.Case {
	path := pathFor(i, 2)
	s := b.Subject(path, b.Any(schedules))
	e := b.Pick(milestones, 5)
	offsets := b.Rng.Perm(40)[:5]
	for k := range offsets {
		offsets[k] += 3
	}
	sort.Ints(offsets)
	// Two milestones are stated relative to one of the three absolutely dated
	// ones. (A slice, not map iteration, keeps the draws deterministic.)
	order := b.Rng.Perm(5)
	absolute := order[2:]
	relative := map[int]int{}
	for _, k := range order[:2] {
		relative[k] = absolute[b.Rng.Intn(len(absolute))]
	}
	ids := make([]string, 5)
	for _, k := range b.Rng.Perm(5) {
		var line string
		if r, rel := relative[k]; rel {
			gap, dir := offsets[k]-offsets[r], "after"
			if gap < 0 {
				gap, dir = -gap, "before"
			}
			line = b.Phrase(orderLines, "rel", "s", s.Name, "e", e[k], "k", fmt.Sprint(gap), "dir", dir, "r", e[r])
		} else {
			line = b.Phrase(orderLines, "abs", "s", s.Name, "e", e[k], "d", b.day(offsets[k]))
		}
		ids[k] = b.Record(s, b.Rng.Intn(Waves), line, e...)
		b.Filler()
	}
	t := 2 + b.Rng.Intn(3)
	gold := e[t-1]
	return valueCase("event-order", path, i, b.Phrase(orderLines, "ask", "a", s.Ask, "t", e[t]), gold, without(e, gold, e[t]), ids...)
}
