package families

import (
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
)

// Indirection families: the answer is plainly stated, but reaching it needs a
// join the question does not spell out: a date inside an interval, a stated
// day that beats the filing order, a nickname, an initial, or a firm behind a
// person.

var asOfLines = Grammar{
	"from": {"From {d}, {s} is in {v}.", "Starting {d}, {s} runs out of {v}."},
	"move": {"From {d}, {s} moves to {v}.", "As of {d}, {s} switches to {v}."},
	"ask": {
		"Where was {a} as of {q}? Just the city.",
		"Which city was {a} in on {q}?",
	},
}

// AsOf dates three locations of an event and asks for the one that held on a
// day strictly inside the middle interval, at least two days after the
// change.
func AsOf(b *Builder, i int) challenge.Case {
	path := pathFor(i, 0)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 3)
	d0 := 1 + b.Rng.Intn(10)
	d1 := d0 + 6 + b.Rng.Intn(10)
	d2 := d1 + 6 + b.Rng.Intn(10)
	q := d1 + 2 + b.Rng.Intn(d2-d1-2)
	w := b.Spread(3)
	b.Record(s, w[0], b.Phrase(asOfLines, "from", "s", s.Name, "d", b.day(d0), "v", c[0]), c[0])
	b.Filler()
	ev := b.Record(s, w[1], b.Phrase(asOfLines, "move", "s", s.Name, "d", b.day(d1), "v", c[1]), c[1])
	b.Filler()
	next := b.Record(s, w[2], b.Phrase(asOfLines, "move", "s", s.Name, "d", b.day(d2), "v", c[2]), c[2])
	return valueCase("as-of", path, i, b.Phrase(asOfLines, "ask", "a", s.Ask, "q", b.dayFixed(q)), c[1], []string{c[0], c[2]}, ev, next)
}

var forwardLines = Grammar{
	"set":  {"{s} will be in {v}.", "{s} is booked in {v}."},
	"move": {"As of {d}, {s} moved to {v}.", "Since {d}, {s} has been in {v}."},
	"old": {
		"Forwarding an old message from {d} that just resurfaced: \"{s} will be in {v}.\"",
		"Found this note from {d} in my drafts: \"{s} is booked in {v}.\"",
	},
	"olddated": {
		"Forwarding a message dated {d}: \"{s} will be in {v}.\"",
	},
}

// ForwardedOld states a move with its date, then, in the last wave, forwards
// an older message that still names the earlier city. The forwarded note is
// the newest filing but the oldest statement: the move still stands. An agent
// that orders by filing time takes the forwarded city.
func ForwardedOld(b *Builder, i int) challenge.Case {
	path := pathFor(i, 1)
	s := b.Subject(path, b.Any(events))
	c := b.Pick(cities, 2)
	d0 := 1 + b.Rng.Intn(8)
	d1 := d0 + 5 + b.Rng.Intn(10)
	w := b.Spread(2)
	b.Record(s, w[0], b.Phrase(forwardLines, "set", "s", s.Name, "v", c[0]), c[0])
	b.Filler()
	ev := b.Record(s, w[1], b.Phrase(forwardLines, "move", "s", s.Name, "d", b.day(d1), "v", c[1]), c[1])
	variant := "old"
	if i%2 == 1 {
		variant = "olddated"
	}
	b.Record(s, Waves-1, b.Phrase(forwardLines, variant, "s", s.Name, "d", b.day(d0), "v", c[0]), c[0])
	return valueCase("forwarded-"+variant, path, i, b.Phrase(cityLines, "ask", "a", s.Ask), c[1], []string{c[0]}, ev)
}

var nickLines = Grammar{
	"quoted": {
		"My colleague {p}, everyone calls them \"{n}\", is joining the {t} team.",
		"Meet {p} from the {t} team; they go by \"{n}\".",
	},
	"bare": {
		"{p} goes by {n} around the {t} team.",
	},
	"old":    {"{p}'s work email is {e}.", "For work, reach {p} at {e}."},
	"change": {"{n} has a new work email: {e}. The old one stops working Friday.", "Heads up, {n} switched to {e} for work."},
	"other":  {"{p}'s work email is {e}.", "Work contact for {p}: {e}."},
	"ask":    {"What is {p}'s work email now?", "Which work address should I use for {p} today?"},
}

// NicknameEmail introduces a person with a nickname, states their work email
// under their full name, and changes it in a note that uses only the
// nickname. The question uses the full name. The bare variant introduces the
// nickname without quotes.
func NicknameEmail(b *Builder, i int) challenge.Case {
	people := b.People(2)
	p, other := people[0], people[1]
	nick := b.Pick(nicknames, 1)[0]
	team := b.Pick(projects, 1)[0]
	d := b.Pick(domains, 2)
	e0, e1, e2 := email(p, d[0]), strings.ToLower(nick+"."+p[strings.LastIndex(p, " ")+1:])+"@"+d[1], email(other, d[0])
	// Only the quoted form splits the frontier panel (8/8, 4/8, 4/8); the bare
	// one was solved by every agent.
	variant := "quoted"
	w := b.Spread(4)
	intro := b.Note(w[0], b.Phrase(nickLines, variant, "p", p, "n", nick, "t", team), p, nick)
	b.Note(w[1], b.Phrase(nickLines, "old", "p", p, "e", e0), p, e0)
	b.Note(w[2], b.Phrase(nickLines, "other", "p", other, "e", e2), other, e2)
	b.Filler()
	ev := b.Note(w[3], b.Phrase(nickLines, "change", "n", nick, "e", e1), nick, e1)
	return valueCase("nickname-email-"+variant, PathRecords, i, b.Phrase(nickLines, "ask", "p", p), e1, []string{e0, e2}, intro, ev)
}

var initialLines = Grammar{
	"directory": {"Directory for the {t} team: {p1} ({e1}), {p2} ({e2}) and {r} ({e3}).", "Contacts, {t} team: {p1} at {e1}; {p2} at {e2}; {r} at {e3}."},
	"holder":    {"{r} is the person holding {s}.", "{s} sits with {r} for now."},
	"handover":  {"{i}. {l} took over {s} from {r}.", "{s} moved from {r} to {i}. {l}."},
	"ask": {
		"Who holds {a} now? Give me their email address.",
		"After every recorded reassignment, what email reaches the person holding {a}?",
	},
}

// InitialsContact hands a thread to someone named only by initial and
// surname. A directory lists them, a namesake with another initial, and the
// previous holder; the question wants the new holder's address.
func InitialsContact(b *Builder, i int) challenge.Case {
	path := pathFor(i, 2)
	s := b.Subject(path, b.Any(threads))
	team := b.Pick(projects, 1)[0]
	surname := b.Any(surnames)
	var p1, p2 string
	for p1 == "" || p2 == "" || p1[0] == p2[0] || b.used[p1] || b.used[p2] {
		p1, p2 = b.Any(firstNames)+" "+surname, b.Any(firstNames)+" "+surname
	}
	b.used[p1], b.used[p2] = true, true
	r := b.Person()
	d := b.Any(domains)
	e1, e2, e3 := email(p1, d), email(p2, d), email(r, d)
	w := b.Spread(3)
	dir := b.Record(s, w[0], b.Phrase(initialLines, "directory", "t", team, "p1", p1, "e1", e1, "p2", p2, "e2", e2, "r", r, "e3", e3), p1, p2, r, e1, e2, e3)
	b.Record(s, w[1], b.Phrase(initialLines, "holder", "s", s.Name, "r", r), r)
	b.Filler()
	ev := b.Record(s, w[2], b.Phrase(initialLines, "handover", "s", s.Name, "r", r, "i", p1[:1], "l", surname), r, surname)
	b.Context(s, 1)
	question, variant := initialLines["ask"][0], "holds"
	if i%2 == 1 {
		question, variant = initialLines["ask"][1], "every"
	}
	return valueCase("initials-contact-"+variant, path, i, capFirst(fill(question, "a", s.Ask)), e1, []string{e2, e3}, dir, ev)
}

var firmLines = Grammar{
	"first": {"{s}: the catering contract went to {f}.", "Catering for {s} is with {f}."},
	"move":  {"{s}: we moved the catering contract to {g}; {p} is their coordinator.", "Catering for {s} now sits with {g}, and {p} will be our contact there."},
	"ask": {
		"Which firm holds the catering contract for {a} now?",
		"Who holds the catering contract for {a} now?",
	},
}

// OwnerFirm moves a contract from one firm to another and names the new
// firm's coordinator in the same note. The question asks who holds the
// contract: the firm, not the person.
func OwnerFirm(b *Builder, i int) challenge.Case {
	path := pathFor(i, 0)
	s := b.Subject(path, b.Any(events))
	f := b.Pick(firms, 2)
	p := b.Person()
	w := b.Spread(2)
	b.Record(s, w[0], b.Phrase(firmLines, "first", "s", s.Name, "f", f[0]), f[0])
	b.Filler()
	ev := b.Record(s, w[1], b.Phrase(firmLines, "move", "s", s.Name, "g", f[1], "p", p), f[1], p)
	b.Context(s, 1)
	// "Which firm" was solved by every agent; "who" makes the DrPresident line
	// answer the coordinator (2/4) while the Hades line names the firm (4/4).
	ask, variant := firmLines["ask"][1], "who"
	return valueCase("owner-firm-"+variant, path, i, fill(ask, "a", s.Ask), f[1], []string{f[0], p}, ev)
}
