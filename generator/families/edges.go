package families

import (
	"fmt"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// Edge families (v8). Same two-item list shape as lists.go: one item is told
// under the person's full name, the other in a way a careful reader still
// follows but that sits just past a limit many current agents share (notes/
// 2026-11-16_edge_tracks.md). Look-alike items of other people are the
// distractors, so naming them, or declining, zeroes the case.

var edgeLines = Grammar{
	"fav": {
		"{p}'s favourite of all the {c} is {v}, apparently.",
		"Of the {c} round here, {p} rates {v} the highest.",
	},
	"fix": {
		"Correction to what I told you about {p}: it isn't {x}, I mixed that up. The one {p} loves is {v}.",
		"Ignore my earlier note about {p} and {x}, I got it wrong. It's {v} that {p} swears by.",
		"Scratch what I said about {p} and {x}. I muddled two people; {p}'s place is {v}.",
	},
	"member": {
		"{p} gave me their membership number, {num}, so I can book things for them.",
		"For bookings on {p}'s behalf, the membership number is {num}.",
	},
	"receipt": {
		"Membership {num} shows up on the receipts from {v} every week, one of the {c} they never skip.",
		"The receipts for membership {num} all come from {v}, a regular among the {c}.",
	},
	"ack": {"Noted.", "Got it.", "Okay, saved that.", "Thanks, noted."},
}

// leads open a conversation with someone by full name; the later turn of the
// same conversation speaks of them only by the phrase.
var leads = [][2]string{
	{"Had a long lunch with {p} today, we talked for ages.", "my lunch companion"},
	{"Shared a cab back from the airport with {p}.", "the person I shared the cab with"},
	{"Sat next to {p} on the train this morning.", "the person I sat next to on the train"},
	{"Spent the evening helping {p} put together a wardrobe.", "the friend I helped with the wardrobe"},
	{"Got paired with {p} at the climbing wall tonight.", "my climbing partner"},
	{"Waited with {p} at the vet's for over an hour.", "the person I waited with at the vet's"},
	{"Was on the same coach tour as {p} all day.", "my seatmate on the coach tour"},
	{"Volunteered at the food bank with {p} this afternoon.", "the volunteer I worked with"},
	{"{p} drove me back from the station when the last train was cancelled.", "the one who drove me back from the station"},
	{"Ran the parents' raffle stall with {p} today.", "the parent I ran the raffle stall with"},
}

// Name shapes the agents' two-capitalised-words name pattern misreads. Each
// gives the asked person and a different, real person whose name the pattern
// confuses with it.
var doubleFirsts = [][2]string{{"Mary", "Ann"}, {"Anna", "Maria"}, {"Jean", "Luc"}, {"Rose", "Marie"}, {"Lou", "Ellen"}, {"Ella", "Mae"}}
var particles = []string{"van der", "de la", "van", "von", "da", "de", "di", "le"}
var particleStems = []string{"Berg", "Heide", "Linden", "Rosa", "Silva", "Costa", "Graaf", "Meer", "Vries", "Bosch", "Luca", "Fontaine"}
var apostropheStems = []string{"Neill", "Connor", "Brien", "Hara", "Donnell", "Reilly", "Sullivan", "Dowd"}
var macStems = []string{"Allister", "Kenna", "Bride", "Gregor", "Intyre", "Kinnon", "Leod", "Farlane"}

// nameShape returns [asked, confusable] for the i-th shape.
func (b *Builder) nameShape(i int) [2]string {
	for {
		first, sur := b.Any(firstNames), b.Any(surnames)
		var out [2]string
		switch i % 5 {
		case 0: // "Mary Ann Delgado" reads as "Mary Ann" or "Ann Delgado"
			d := doubleFirsts[b.Rng.Intn(len(doubleFirsts))]
			out = [2]string{d[0] + " " + d[1] + " " + sur, d[1] + " " + sur}
		case 1: // "Lars van der Berg": no two-capital run at all
			stem := b.Any(particleStems)
			out = [2]string{first + " " + b.Any(particles) + " " + stem, first + " " + stem}
		case 2: // "Siobhan O'Neill"
			stem := b.Any(apostropheStems)
			out = [2]string{first + " O'" + stem, first + " " + stem}
		case 3: // "Siobhan McAllister"
			stem := b.Any(macStems)
			out = [2]string{first + " Mc" + stem, first + " " + stem}
		default: // "Anna Smith-Jones" reads as "Anna Smith"
			other := b.Any(surnames)
			if other == sur {
				continue
			}
			out = [2]string{first + " " + sur + "-" + other, first + " " + sur}
		}
		if !b.used[out[0]] && !b.used[out[1]] {
			b.used[out[0]], b.used[out[1]] = true, true
			return out
		}
	}
}

// NameShapeList: both shops are told under a name outside the two-word shape
// ("Mary Ann Delgado", "Lars van der Berg", "Siobhan O'Neill", "Anna Smith-
// Jones"), the second in a record that never names the shop kind; a different
// person whose name the pattern confuses with it ("Ann Delgado", "Lars Berg",
// "Anna Smith") has two shops of their own.
func NameShapeList(b *Builder, i int) challenge.Case {
	people := b.nameShape(i)
	l := b.newList(people, people)
	l.items("plain")
	return l.finish("name-shape-list")
}

// CorrectionList: one shop under the full name; the other is first given
// wrongly and then corrected in a later record that says the first was a
// mix-up. The wrong one is a distractor, as is a namesake's shop: an agent
// that presents both states as a disagreement, or keeps the first, scores 0.
func CorrectionList(b *Builder, i int) challenge.Case {
	people := b.namesakes()
	l := b.newList(people, people)
	wrong := l.shops[3]
	early, late := l.w[2], l.w[6]
	if early > late {
		early, late = late, early
	}
	l.ids = append(l.ids,
		l.say(1, "full", "p", people[0], "v", l.shops[0]),
		l.b.Note(early, capFirst(b.Phrase(edgeLines, "fav", "p", people[0], "v", wrong, "c", l.kind[0])), l.protect...),
		l.b.Note(late, capFirst(b.Phrase(edgeLines, "fix", "p", people[0], "x", wrong, "v", l.shops[1])), l.protect...))
	l.say(4, "full", "p", people[1], "v", l.shops[2])
	return l.finish("correction-list")
}

// SessionRelay: one shop under the full name; the other only in a later turn
// of the conversation that named the person, which speaks of them by a phrase
// ("My lunch companion swears by ..."). Another conversation is seeded between
// the two turns, so they are 40 to 160 minutes apart with a record of another
// session in between. Two look-alike conversations about other people use
// other leads, with their own shops as distractors.
func SessionRelay(b *Builder, i int) challenge.Case {
	people := [2]string{b.Person(), b.Person()}
	l := b.newList(people, people)
	order := b.Rng.Perm(len(leads))
	relay := func(k int, person string, lead [2]string, shop string) (string, string) {
		w := l.w[k]
		s := b.Session(w)
		intro := b.SayIn(w, s, fill(lead[0], "p", person), b.Any(acks))
		line := filler[b.Rng.Intn(len(filler))]
		b.Say(w, line[0], line[1])
		turn := b.SayIn(w, s, capFirst(b.Phrase(listLines, "plain", "n", lead[1], "v", shop)), b.Any(edgeLines["ack"]))
		return intro, turn
	}
	intro, turn := relay(2, people[0], leads[order[0]], l.shops[1])
	l.ids = append(l.ids, l.say(1, "full", "p", people[0], "v", l.shops[0]), intro, turn)
	relay(4, people[1], leads[order[1]], l.shops[2])
	other := b.Person()
	relay(5, other, leads[order[2]], l.shops[3])
	return l.finish("session-relay")
}

// NumberList: one shop under the full name; the other only on receipts filed
// under a four-digit membership number that another record gives as the
// person's. A receipt under the same digits in another order names a
// distractor shop, and a namesake's shop is the other.
func NumberList(b *Builder, i int) challenge.Case {
	people := b.namesakes()
	l := b.newList(people, people)
	num, swapped := b.fourDigits()
	if i%2 == 1 {
		num, swapped = "#"+num, "#"+swapped
	}
	l.ids = append(l.ids,
		l.say(1, "full", "p", people[0], "v", l.shops[0]),
		l.b.Note(l.w[2], capFirst(b.Phrase(edgeLines, "member", "p", people[0], "num", num)), append(l.protect, num)...),
		l.b.Note(l.w[5], capFirst(b.Phrase(edgeLines, "receipt", "num", num, "v", l.shops[1], "c", l.kind[0])), append(l.protect, num)...))
	l.b.Note(l.w[6], capFirst(b.Phrase(edgeLines, "receipt", "num", swapped, "v", l.shops[3], "c", l.kind[0])), append(l.protect, swapped)...)
	l.say(4, "full", "p", people[1], "v", l.shops[2])
	return l.finish("number-list")
}

// fourDigits returns a fresh four-digit number and the same digits with two
// neighbours swapped (a different number).
func (b *Builder) fourDigits() (string, string) {
	for {
		n := fmt.Sprintf("%04d", 1000+b.Rng.Intn(9000))
		k := b.Rng.Intn(3)
		s := []byte(n)
		s[k], s[k+1] = s[k+1], s[k]
		sw := string(s)
		if sw != n && sw[0] != '0' && !b.used[n] && !b.used[sw] {
			b.used[n], b.used[sw] = true, true
			return n, sw
		}
	}
}

// RoleNoteList: RoleList with each role binding told inside a long note about
// something else, so the record that ties the role to the person is well over
// the length the agents' link step reads.
func RoleNoteList(b *Builder, i int) challenge.Case {
	people := b.namesakes()
	teams := b.Pick(projects, 2)
	l := b.newList(people, [2]string{"the " + teams[0] + " lead", "the " + teams[1] + " lead"})
	l.ids = append(l.ids, l.sayInNote(0, b.Any(events), "role", "p", people[0], "t", teams[0]))
	l.sayInNote(3, b.Any(events), "role", "p", people[1], "t", teams[1])
	l.items("ref")
	return l.finish("role-note-list")
}

// chatter names a person and says nothing a question asks, so a person's own
// records are many and the ones that matter must be read for, not counted.
var chatter = []string{
	"{p} rang this evening, mostly to complain about the weather.",
	"{p} sent me a voice message I still have to listen to.",
	"{p} has a report due on Friday and is barely sleeping.",
	"{p} oiled the hinges on the back door at last, no more creaking.",
	"{p} spent twenty minutes describing the upstairs neighbour's parrot.",
	"{p} left an umbrella at mine, the second one this month.",
	"{p} can't decide between grey and sage for the kitchen walls.",
	"{p} is off sick with a sore throat and a bad cough.",
	"{p} wanted to know if a rattling exhaust is worth worrying about.",
	"{p} and I spent ages trying to assemble a flat-pack shelf.",
	"{p} has started getting up early to go for a run.",
	"{p} is planning to grow tomatoes on the balcony.",
}

// Chatter seeds n chatter records about person in random waves, and a plain
// filler pair beside every other one. Generation time grows with the pair
// count and a run is capped at 60 s, so an 80-case pack stays near 1,200 pairs.
func (b *Builder) Chatter(person string, n int) {
	for k := 0; k < n; k++ {
		b.Say(b.Rng.Intn(Waves), capFirst(fill(b.Any(chatter), "p", person)), b.Any(acks))
		if k%2 == 0 {
			b.Filler()
		}
	}
}

// joinShops is the gold answer text.
func joinShops(v []string) string { return strings.Join(v, ", ") }

var _ = protocol.AnswerList
