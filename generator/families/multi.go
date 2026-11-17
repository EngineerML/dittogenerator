package families

import (
	"fmt"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// Multi-link lists (v8.1). A two-item list is solved (score >= 0.5) by anyone
// who finds the full-name item, so most cases come out trivial. Here a case asks
// for three things: one under the full name and one through each of two
// different edge links. An agent needs the full-name item and at least one link
// to reach 2/3; agents fail different links, so cases split the panel. Each link
// brings a distractor of its own (another person's relay, a number with two
// digits swapped, a relative under the other title, a place given up), and
// naming one zeroes the case.

// multi is one case: the person, the shop kind, the shops drawn for it, and the
// evidence and distractors its links add.
type multi struct {
	b       *Builder
	kind    [2]string
	person  string
	first   string
	surname string
	pronoun [2]string // she/her or he/his
	shops   []string
	next    int
	items   []string
	traps   []string
	ids     []string
	protect []string
	waves   []int
	wnext   int
}

var pronouns = [][2]string{{"she", "her"}, {"he", "his"}}

func (b *Builder) newMulti(person, first, surname string) *multi {
	kind := listKinds[b.Len()%len(listKinds)]
	var shops []string
	for _, i := range b.Rng.Perm(len(shopNames)) {
		if len(shops) == 5 {
			break
		}
		if s := shopNames[i] + " " + kind[1]; !b.used[s] {
			b.used[s] = true
			shops = append(shops, s)
		}
	}
	if len(shops) < 5 {
		panic("value pool exhausted: shop names for " + kind[0])
	}
	m := &multi{b: b, kind: kind, person: person, first: first, surname: surname, pronoun: pronouns[b.Rng.Intn(2)], shops: shops, waves: b.Rng.Perm(Waves)}
	m.protect = append([]string{person, first, surname, "membership", "receipts"}, shops...)
	return m
}

func (m *multi) shop() string { s := m.shops[m.next]; m.next++; return s }
func (m *multi) wave() int    { w := m.waves[m.wnext%Waves]; m.wnext++; return w }

// two returns two of the case's waves in increasing order.
func (m *multi) two() (int, int) {
	a, c := m.wave(), m.wave()
	if a > c {
		a, c = c, a
	}
	return a, c
}

func (m *multi) note(w int, line string, extra ...string) string {
	return m.b.Note(w, capFirst(line), append(m.protect, extra...)...)
}

var multiLines = Grammar{
	"full": {
		"{p} loves {v} and goes there most weekends; {s} says none of the other {c} come close.",
		"{p} keeps telling me {v} is the best of the {c} around here, {s} goes every week.",
	},
	"title": {
		"{t} {sur} is at {v} every week without fail.",
		"Bumped into {t} {sur} at {v} again, a regular there.",
	},
	"used": {
		"{p} goes to {v} every week, has done for ages.",
		"{p} has a loyalty card at {v} and uses it all the time.",
	},
	"dropped": {
		"{f} has had enough of {v} after the last price rise and isn't going back.",
		"{f} let the card for {v} lapse and won't be renewing it.",
		"{f} fell out with the owner of {v} over a refund and is done with the place.",
	},
	"found": {
		"{f} found {v} instead and now can't stop talking about it.",
		"These days {f} goes to {v}, a lovely find.",
	},
	// soft forms (v8.3): the person is named in full, the relative's tie is said.
	"droppedFull": {
		"{p} stopped going to {v} after the last price rise and isn't going back.",
		"{p} let the card for {v} lapse and won't be renewing it.",
		"{p} has given up on {v} for good after a row with the staff.",
	},
	"foundFull": {
		"{p} goes to {v} instead now and can't stop talking about it.",
		"Since leaving the old place, {p} has been going to {v}.",
	},
	"titleRel": {
		"{p}'s {rel}, {t} {sur}, is a regular at {v}.",
		"{t} {sur}, {p}'s {rel}, is at {v} every week without fail.",
	},
	"ask": {
		"Which {c} does {p} go to these days?",
		"Which two {c} does {p} love?",
		"What are the {c} {p} keeps going back to?",
	},
}

// linkSession: the item is told only in the next turn of a conversation that
// named the person, an hour or two later with another conversation in between;
// another person's conversation under a different lead names the distractor.
func (m *multi) linkSession() {
	item, trap := m.shop(), m.shop()
	order := m.b.Rng.Perm(len(leads))
	relay := func(person string, lead [2]string, shop string) (string, string) {
		w := m.wave()
		s := m.b.Session(w)
		intro := m.b.SayIn(w, s, fill(lead[0], "p", person), m.b.Any(acks))
		line := filler[m.b.Rng.Intn(len(filler))]
		m.b.Say(w, line[0], line[1])
		turn := m.b.SayIn(w, s, capFirst(m.b.Phrase(listLines, "plain", "n", lead[1], "v", shop)), m.b.Any(edgeLines["ack"]))
		return intro, turn
	}
	a, c := relay(m.person, leads[order[0]], item)
	relay(m.b.Person(), leads[order[1]], trap)
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a, c)
}

// linkNumber: the item is on receipts under a four-digit membership number
// another record gives as the person's; the same digits in another order are
// on a receipt for the distractor.
func (m *multi) linkNumber() {
	item, trap := m.shop(), m.shop()
	num, swapped := m.b.fourDigits()
	if m.b.Rng.Intn(2) == 1 {
		num, swapped = "#"+num, "#"+swapped
	}
	a := m.note(m.wave(), m.b.Phrase(edgeLines, "member", "p", m.person, "num", num), num)
	c := m.note(m.wave(), m.b.Phrase(edgeLines, "receipt", "num", num, "v", item, "c", m.kind[0]), num)
	m.note(m.wave(), m.b.Phrase(edgeLines, "receipt", "num", swapped, "v", trap, "c", m.kind[0]), swapped)
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a, c)
}

// linkTitle: the item is told under the title that fits the person (the
// full-name record says she or he); a relative under the other title has the
// distractor.
func (m *multi) linkTitle() {
	item, trap := m.shop(), m.shop()
	mine, other := "Ms.", "Mr."
	if m.pronoun[0] == "he" {
		mine, other = "Mr.", "Ms."
	}
	a := m.note(m.wave(), m.b.Phrase(multiLines, "title", "t", mine, "sur", m.surname, "v", item))
	m.note(m.wave(), m.b.Phrase(multiLines, "title", "t", other, "sur", m.surname, "v", trap))
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a)
}

// linkGivenUp: the person used a place, later gave it up in loose wording
// (first name only) and found another; the given-up place is the distractor.
func (m *multi) linkGivenUp(soft bool) {
	item, trap := m.shop(), m.shop()
	early, late := m.two()
	m.note(early, m.b.Phrase(multiLines, "used", "p", m.person, "v", trap))
	dropped, found := m.b.Phrase(multiLines, "dropped", "f", m.first, "v", trap), m.b.Phrase(multiLines, "found", "f", m.first, "v", item)
	if soft {
		dropped, found = m.b.Phrase(multiLines, "droppedFull", "p", m.person, "v", trap), m.b.Phrase(multiLines, "foundFull", "p", m.person, "v", item)
	}
	a := m.note(late, dropped)
	c := m.note(late, found)
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a, c)
}

// linkRoleNote: the item is told under a role whose holder is named inside a
// long note about something else; another team's lead has the distractor.
func (m *multi) linkRoleNote() {
	item, trap := m.shop(), m.shop()
	teams := m.b.Pick(projects, 2)
	other := m.b.Person()
	a := m.b.LongNote(m.wave(), m.b.Any(events), capFirst(m.b.Phrase(listLines, "role", "p", m.person, "t", teams[0])), m.protect...)
	m.b.LongNote(m.wave(), m.b.Any(events), capFirst(m.b.Phrase(listLines, "role", "p", other, "t", teams[1])), m.protect...)
	c := m.note(m.wave(), m.b.Phrase(listLines, "plain", "n", "the "+teams[0]+" lead", "v", item))
	m.note(m.wave(), m.b.Phrase(listLines, "plain", "n", "the "+teams[1]+" lead", "v", trap))
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a, c)
}

// linkShape: the person's name is outside the two-word shape; the item is
// told under it in a record that never names the shop kind, and a person whose
// name the two-word pattern confuses with it has the distractor. Only for
// name-shape people (multi.confusable set).
func (m *multi) linkShape(confusable string) {
	item, trap := m.shop(), m.shop()
	a := m.note(m.wave(), m.b.Phrase(listLines, "plain", "n", m.person, "v", item))
	m.note(m.wave(), m.b.Phrase(listLines, "full", "p", confusable, "v", trap, "c", m.kind[0]))
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, a)
}

var relatives = map[string][]string{"Mr.": {"husband", "brother", "father", "son"}, "Ms.": {"wife", "sister", "mother", "daughter"}}

// trapTitle: a relative under the title that does not fit the person (the
// full-name record says she or he) has a place of their own. No item. Soft:
// the record names the tie to the person.
func (m *multi) trapTitle(soft bool) {
	trap := m.shop()
	other := "Mr."
	if m.pronoun[0] == "he" {
		other = "Ms."
	}
	if soft {
		rel := relatives[other][m.b.Rng.Intn(len(relatives[other]))]
		m.note(m.wave(), m.b.Phrase(multiLines, "titleRel", "p", m.person, "rel", rel, "t", other, "sur", m.surname, "v", trap))
	} else {
		m.note(m.wave(), m.b.Phrase(multiLines, "title", "t", other, "sur", m.surname, "v", trap))
	}
	m.traps = append(m.traps, trap)
}

// trapGivenUp: the person used a place and later gave it up in loose wording
// under the first name only. No item. Soft: the giving up names the person in
// full.
func (m *multi) trapGivenUp(soft bool) {
	trap := m.shop()
	early, late := m.two()
	m.note(early, m.b.Phrase(multiLines, "used", "p", m.person, "v", trap))
	if soft {
		m.note(late, m.b.Phrase(multiLines, "droppedFull", "p", m.person, "v", trap))
	} else {
		m.note(late, m.b.Phrase(multiLines, "dropped", "f", m.first, "v", trap))
	}
	m.traps = append(m.traps, trap)
}

// multiPlans (v8.3) pair one item link (the second item; "shape" makes the
// person a name-shape person; "found" is the given-up link, whose replaced
// place is its own trap) with one trap ("title": a relative under the other
// title; "givenup": a place the person gave up; "-": none) at a strength.
// With two items, a case is solved by anyone who finds the full-name item and
// names no trap, so the trap decides solved / failed and the link item only
// lowers the mean. v8.2 had every trap at full strength: on the released-agent
// panel each agent named one in 34-71% of cases and 15 of 35 cases went
// unsolved. Live round columns that rank keep unsolved at 0-11 of 70, so v8.3
// writes about half its traps soft (the person named in full, the relative's
// tie said) as a ladder: hard cases split strong panels, soft ones weak.
//
// v8.4: on the same panel v8.3 scored about 0.08 (15 discriminating, 9
// unsolved, 11 trivial). The unsolved cases sat in title-givenup-hard (0-2
// solves of 3 per agent) and shape-givenup-hard (0-1); they become
// title-givenup-soft and shape-title-hard.
//
// v8.5: v8.4 still left 10 of 35 unsolved, mostly where a given-up trap sat on
// top of a link's own trap (number, shape, title + givenup: 0-2 solves of 3-4
// per agent). Every link already brings a trap of its own, so the given-up
// trap stays only with the session link (solved by every agent) and the
// given-up link itself; the rest carry a title trap or no extra one.
var multiPlans = [][3]string{
	{"number", "title", "hard"}, {"session", "title", "hard"}, {"shape", "title", "soft"},
	{"number", "-", "soft"}, {"session", "givenup", "soft"}, {"shape", "-", "soft"},
	{"title", "-", "soft"}, {"number", "title", "soft"}, {"session", "title", "soft"},
	{"shape", "title", "soft"}, {"found", "-", "soft"},
}

// MultiList builds the i-th rotation plan's case.
func MultiList(b *Builder, i int) challenge.Case {
	plan := multiPlans[b.Len()%len(multiPlans)]
	soft := plan[2] == "soft"
	var person, first, surname, confusable string
	if plan[0] == "shape" {
		shape := b.nameShape(b.Len())
		person, confusable = shape[0], shape[1]
		parts := strings.Fields(person)
		first, surname = parts[0], parts[len(parts)-1]
	} else {
		person = b.OwnPerson()
		parts := strings.Fields(person)
		first, surname = parts[0], parts[1]
	}
	m := b.newMulti(person, first, surname)
	full := m.shop()
	m.items = append(m.items, full)
	m.ids = append(m.ids, m.note(m.wave(), m.b.Phrase(multiLines, "full", "p", person, "v", full, "c", m.kind[0], "s", m.pronoun[0])))
	switch plan[0] {
	case "shape":
		m.linkShape(confusable)
	case "session":
		m.linkSession()
	case "number":
		m.linkNumber()
	case "title":
		m.linkTitle()
	case "role":
		m.linkRoleNote()
	case "found":
		m.linkGivenUp(soft)
	}
	switch plan[1] {
	case "title":
		if plan[0] != "title" {
			m.trapTitle(soft)
		}
	case "givenup":
		m.trapGivenUp(soft)
	}
	b.Chatter(person, 3)
	kind := "multi-" + plan[0] + "-" + plan[1] + "-" + plan[2]
	return challenge.Case{
		MemoryCase: protocol.MemoryCase{
			ID:                fmt.Sprintf("%s-%d", kind, b.Len()),
			QuestionType:      kind,
			Question:          b.Phrase(multiLines, "ask", "c", m.kind[0], "p", person),
			ExpectedAnswer:    strings.Join(m.items, ", "),
			AnswerKind:        protocol.AnswerList,
			AnswerItems:       m.items,
			DistractorAnswers: m.traps,
		},
		EvidencePairIDs: m.ids,
	}
}
