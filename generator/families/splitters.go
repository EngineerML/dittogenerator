package families

import (
	"github.com/ditto-assistant/dittobench-datagen/challenge"
)

// Splitter families: ordinary memory questions worded the way people ask
// them, where one line of the top agents routes the wording somewhere else
// before it reads memory ("what we decided" reads as an instruction to act,
// "the up-to-date status" as a live lookup, a misspelt "emial" misses a
// literal contact pattern). Each answer is plainly in memory; the case only
// tests whether the agent reads it. They are a quarter of a pack at most,
// since one regex change in a new agent version can retire them.

var decidedLines = Grammar{
	"weigh":   {"For {s} we're weighing {v} and {w}.", "Shortlist for {s}: {v} or {w}."},
	"suggest": {"{p} also suggested {v} for {s}.", "{p} pitched {v} for {s} too."},
	"decide":  {"Final call on {s}: we agreed on {v}.", "Settled it, {s} goes to {v}."},
	"ask": {
		"Remind me what we decided for {a}, which city?",
		"Which city did we settle on for {a}?",
	},
}

// Decided weighs two cities, adds a third suggestion, then records the
// decision.
func Decided(b *Builder, i int) challenge.Case {
	path := pathFor(i, 1)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 3)
	p := b.Person()
	w := b.Spread(3)
	b.Record(s, w[0], b.Phrase(decidedLines, "weigh", "s", s.Name, "v", c[0], "w", c[1]), c[0], c[1])
	b.Record(s, w[1], b.Phrase(decidedLines, "suggest", "s", s.Name, "v", c[2], "p", p), c[2], p)
	b.Filler()
	ev := b.Record(s, w[2], b.Phrase(decidedLines, "decide", "s", s.Name, "v", c[1]), c[1])
	ask, variant := decidedLines["ask"][0], "decided"
	if i%2 == 1 {
		ask, variant = decidedLines["ask"][1], "settle"
	}
	return valueCase("decided-"+variant, path, i, fill(ask, "a", s.Ask), c[1], []string{c[0], c[2]}, ev)
}

var statusLines = Grammar{
	"first":  {"{s} is {v}.", "Where {s} stands: {v}."},
	"change": {"{s} is now {v}.", "Status change for {s}: {v}.", "{s} moved to {v}."},
	"ask": {
		"What's the up-to-date status of {a}?",
		"Where does {a} stand at the moment?",
		"What is the status of {a} now?",
	},
}

// LiveStatus moves a project through three multi-word states and asks for
// the current one.
func LiveStatus(b *Builder, i int) challenge.Case {
	path := pathFor(i, 2)
	s := b.Subject(path, b.Any(threads))
	st := b.Pick(statuses, 3)
	w := b.Spread(3)
	b.Record(s, w[0], b.Phrase(statusLines, "first", "s", s.Name, "v", st[0]), st[0])
	b.Filler()
	b.Record(s, w[1], b.Phrase(statusLines, "change", "s", s.Name, "v", st[1]), st[1])
	ev := b.Record(s, w[2], b.Phrase(statusLines, "change", "s", s.Name, "v", st[2]), st[2])
	variants := []string{"uptodate", "stand", "plain"}
	return valueCase("live-status-"+variants[i%3], path, i, fill(statusLines["ask"][i%3], "a", s.Ask), st[2], st[:2], ev)
}

var runnerLines = Grammar{
	"runs":     {"{p} is running {s} for now.", "{p} has {s} in hand."},
	"handover": {"{p} handed {s} over to {q}.", "{q} is running {s} from today; {p} stepped back."},
	"card":     {"{p}'s email: {e}.", "Contact card, {p}: {e}."},
	"typo":     {"whats the emial for whoever runs {a} now?", "who runs {a} these days, and whats there email adress?"},
	"clean":    {"What's the email for whoever runs {a} now?"},
}

// TypoContact hands a project over between two people in long notes and
// keeps each person's email in a short contact card. The question asks for
// the current runner's email, in the typed-in-a-hurry variant with slips a
// literal contact pattern misses, or cleanly as the control. It always takes
// the topic path, where contact cards are added only for contact questions.
func TypoContact(b *Builder, i int) challenge.Case {
	s := b.Subject(PathTopic, "offsite")
	r := b.People(2)
	d := b.Any(domains)
	e0, e1 := email(r[0], d), email(r[1], d)
	w := b.Spread(4)
	b.Record(s, w[0], b.Phrase(runnerLines, "runs", "s", s.Name, "p", r[0]), r[0])
	b.Note(w[1], b.Phrase(runnerLines, "card", "p", r[0], "e", e0), r[0], e0)
	b.Note(w[1], b.Phrase(runnerLines, "card", "p", r[1], "e", e1), r[1], e1)
	b.Record(s, w[2], fill(b.Any(paragraphs)+" {s} is still on the calendar.", "s", capFirst(s.Name)))
	hand := b.Record(s, w[3], b.Phrase(runnerLines, "handover", "s", s.Name, "p", r[0], "q", r[1]), r[0], r[1])
	card := b.Note(w[3], b.Phrase(runnerLines, "card", "p", r[1], "e", e1), r[1], e1)
	variant := "typo"
	if i%3 == 2 {
		variant = "clean"
	}
	return valueCase("typo-contact-"+variant, PathTopic, i, b.Phrase(runnerLines, variant, "a", s.Ask), e1, []string{e0}, hand, card)
}

var handoverLines = Grammar{
	"start": {"{p} is handling {s}.", "{s} is with {p} for now."},
	"cc":    {"Cc'd {d} on {s} for visibility.", "{d} is copied on {s} updates."},
	"take": {
		"{q} took over {s}.",
		"{q} will run point on {s} from here.",
		"{s} is {q}'s now; they have the ball.",
		"{q} owns {s} going forward.",
	},
	"ask": {"Who is handling {a} now?", "Who has {a} at this point?"},
}

// Handover passes a thread between two people with one of several handover
// verbs, each matched by a different agent line's cue list, and copies a
// third person on it as a decoy.
func Handover(b *Builder, i int) challenge.Case {
	path := pathFor(i, 0)
	s := b.Subject(path, b.Any(threads))
	p := b.People(3)
	w := b.Spread(3)
	b.Record(s, w[0], b.Phrase(handoverLines, "start", "s", s.Name, "p", p[0]), p[0])
	b.Record(s, w[1], b.Phrase(handoverLines, "cc", "s", s.Name, "d", p[2]), p[2])
	b.Filler()
	take := handoverLines["take"][i%len(handoverLines["take"])]
	ev := b.Record(s, w[2], capFirst(fill(take, "s", s.Name, "q", p[1])), p[1])
	variants := []string{"took", "point", "ball", "owns"}
	return valueCase("handover-"+variants[i%4], path, i, b.Phrase(handoverLines, "ask", "a", s.Ask), p[1], []string{p[0], p[2]}, ev)
}
