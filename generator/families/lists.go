package families

import (
	"fmt"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// The list families ask for two things one person likes. The first is told
// under the person's full name; the second only through a link the agent has
// to follow (a nickname, initials, a role, a misspelled name), and a
// look-alike person nearby has things of their own. It is the shape behind
// the generator competition's half-credit columns (docs/v14-generator-
// credits.md, SN-352): a list is graded by the share of items named, and a
// distractor or a decline zeroes it. An agent that names only the full-name
// item earns 0.5, which counts as solved; one that hedges about the missing
// second item ("and no record of a second one"), or takes the look-alike's,
// earns 0; one that follows the link earns 1.
//
// Every link is stated or unambiguous, so each case is solvable at full
// credit.
var listLines = Grammar{
	"full": {
		"{p} loves {v} and goes there most weekends.",
		"{p} keeps telling me {v} is the best of the {c} around here.",
		"Ran into {p} at {v} again; it's the favourite of all the {c} for sure.",
	},
	"ref": {
		"{n} swears by {v}, says none of the other {c} come close.",
		"{n} dragged me to {v} today and would not stop praising it.",
		"Apparently {v} is where {n} goes every single week.",
	},
	// "plain" never names the shop kind, so the question's words do not pull
	// the second record up: an agent has to reach it through the link.
	"plain": {
		"{n} swears by {v}, nothing else comes close.",
		"{n} dragged me to {v} today and would not stop praising it.",
		"Apparently {v} is where {n} goes every single week.",
	},
	"nickname": {
		"{p} goes by {n} with pretty much everyone.",
		"Everyone calls {p} {n}, has done for years.",
		"Quick note so I don't mix them up: {p} is the one we call {n}.",
	},
	"role": {
		"{p} took over as the {t} lead this month.",
		"New org chart: {p} is now the {t} lead.",
		"Heads up, the {t} lead is {p} from now on.",
	},
	"ask": {
		// Questions only: an imperative ("Name the two...") reads as a request to
		// act, not a question about memory.
		"Which two {c} does {p} love?",
		"Which two {c} is {p} really fond of?",
		"What are the two {c} {p} keeps going back to?",
	},
}

// listKinds are the shop kinds: the plural used in the text and the noun that
// ends each shop's name. A pack draws them in turn, so no kind runs short.
var listKinds = [][2]string{
	{"bakeries", "Bakehouse"}, {"bookshops", "Books"}, {"coffee roasters", "Roasters"},
	{"delis", "Deli"}, {"florists", "Flowers"}, {"bike shops", "Cycles"},
}

// lookalikes are pairs of real surnames one letter apart, for the misspelled-
// name form: the misspelling is a swap of two letters of the first, the
// look-alike carries the second.
var lookalikes = [][2]string{
	{"Page", "Pace"}, {"Holm", "Holt"}, {"Reid", "Reed"}, {"Kane", "Lane"}, {"Moss", "Ross"},
	{"Hale", "Hall"}, {"Lowe", "Love"}, {"Wade", "Ward"}, {"Kerr", "Carr"}, {"Mann", "Munn"},
	{"Dunn", "Dunne"}, {"Price", "Pryce"}, {"Gray", "Grey"}, {"Marsh", "Nash"}, {"Burke", "Burge"},
	{"Hart", "Hunt"}, {"Cole", "Dole"}, {"Fry", "Foy"}, {"Webb", "Webber"}, {"Lamb", "Lamm"},
}

// listCase is one case of the shape: the person asked about and the look-alike,
// how each is referred to in the second record, the shop kind and four shops.
type listCase struct {
	b       *Builder
	kind    [2]string
	shops   []string
	people  [2]string
	refs    [2]string
	w       []int
	ids     []string
	protect []string
}

func (b *Builder) newList(people [2]string, refs [2]string) *listCase {
	kind := listKinds[b.Len()%len(listKinds)]
	var shops []string
	for _, i := range b.Rng.Perm(len(shopNames)) {
		if len(shops) == 4 {
			break
		}
		if s := shopNames[i] + " " + kind[1]; !b.used[s] {
			b.used[s] = true
			shops = append(shops, s)
		}
	}
	if len(shops) < 4 {
		panic("value pool exhausted: shop names for " + kind[0])
	}
	l := &listCase{b: b, kind: kind, shops: shops, people: people, refs: refs, w: b.Rng.Perm(Waves)}
	l.protect = append(append(append([]string(nil), people[:]...), refs[:]...), shops...)
	return l
}

// say seeds one record in the k-th of the case's waves and returns its pair ID.
func (l *listCase) say(k int, key string, kv ...string) string {
	line := capFirst(l.b.Phrase(listLines, key, append(kv, "c", l.kind[0])...))
	return l.b.Note(l.w[k], line, l.protect...)
}

// sayInNote seeds a line inside a longer note about an unrelated topic, so it
// is stated plainly but does not stand on its own.
func (l *listCase) sayInNote(k int, topic, key string, kv ...string) string {
	line := capFirst(l.b.Phrase(listLines, key, kv...))
	return l.b.LongNote(l.w[k], topic, line, l.protect...)
}

// items seeds both people's two shops: one under the full name, one under the
// reference worded by ref ("ref" or "plain"). The first person's two records
// are the evidence.
func (l *listCase) items(ref string) {
	l.ids = append(l.ids,
		l.say(1, "full", "p", l.people[0], "v", l.shops[0]),
		l.say(2, ref, "n", l.refs[0], "v", l.shops[1]))
	l.say(4, "full", "p", l.people[1], "v", l.shops[2])
	l.say(5, ref, "n", l.refs[1], "v", l.shops[3])
}

func (l *listCase) finish(kind string) challenge.Case {
	gold := l.shops[:2]
	return challenge.Case{
		MemoryCase: protocol.MemoryCase{
			// A family may hold two slots in the rotation, so the round number
			// alone would repeat; the case's place in the pack does not.
			ID:                fmt.Sprintf("%s-%d", kind, l.b.Len()),
			QuestionType:      kind,
			Question:          l.b.Phrase(listLines, "ask", "c", l.kind[0], "p", l.people[0]),
			ExpectedAnswer:    strings.Join(gold, ", "),
			AnswerKind:        protocol.AnswerList,
			AnswerItems:       gold,
			DistractorAnswers: l.shops[2:],
		},
		EvidencePairIDs: l.ids,
	}
}

// namesakes returns two reserved full names that share a first name.
func (b *Builder) namesakes() [2]string {
	first := b.Any(firstNames)
	var out [2]string
	for k := 0; k < 2; {
		if name := first + " " + b.Any(surnames); !b.used[name] {
			b.used[name] = true
			out[k] = name
			k++
		}
	}
	return out
}

// AliasList: the second shop is told under a nickname, and a namesake (same
// first name) has a nickname and shops of their own. Both nicknames are bound
// once in plain words, inside a longer note about something else, and the
// second shop's record never names the shop kind.
func AliasList(b *Builder, i int) challenge.Case {
	people := b.namesakes()
	nick := b.Pick(nicknames, 2)
	l := b.newList(people, [2]string{nick[0], nick[1]})
	l.ids = append(l.ids, l.sayInNote(0, b.Any(events), "nickname", "p", people[0], "n", nick[0]))
	l.sayInNote(3, b.Any(events), "nickname", "p", people[1], "n", nick[1])
	l.items("plain")
	return l.finish("alias-list")
}

// InitialsList: the second shop is told under the initial and surname
// ("E. Page"), and a look-alike shares the surname with another initial
// ("Rosa Page", "R. Page"), so only the initial tells them apart.
func InitialsList(b *Builder, i int) challenge.Case {
	surname := b.Any(surnames)
	var people [2]string
	for people[0] == "" || people[1] == "" || people[0][0] == people[1][0] || b.used[people[0]] || b.used[people[1]] {
		people = [2]string{b.Any(firstNames) + " " + surname, b.Any(firstNames) + " " + surname}
	}
	b.used[people[0]], b.used[people[1]] = true, true
	l := b.newList(people, [2]string{people[0][:1] + ". " + surname, people[1][:1] + ". " + surname})
	l.items("ref")
	return l.finish("initials-list")
}

// RoleList: the second shop is told under a role ("the Brevant lead"), and a
// record elsewhere names who holds it; a namesake holds another team's lead.
func RoleList(b *Builder, i int) challenge.Case {
	people := b.namesakes()
	teams := b.Pick(projects, 2)
	l := b.newList(people, [2]string{"the " + teams[0] + " lead", "the " + teams[1] + " lead"})
	l.ids = append(l.ids, l.say(0, "role", "p", people[0], "t", teams[0]))
	l.say(3, "role", "p", people[1], "t", teams[1])
	l.items("ref")
	return l.finish("role-list")
}

// TypoList: the second shop is told under the person's name with two letters
// of the surname swapped ("Edith Pgae"), in a record that never names the
// shop kind, and a real look-alike one letter away ("Edith Pace") has shops
// of their own.
func TypoList(b *Builder, i int) challenge.Case {
	first := b.Any(firstNames)
	var pair [2]string
	for {
		pair = lookalikes[b.Rng.Intn(len(lookalikes))]
		if !b.used[first+" "+pair[0]] && !b.used[first+" "+pair[1]] {
			break
		}
		first = b.Any(firstNames)
	}
	people := [2]string{first + " " + pair[0], first + " " + pair[1]}
	b.used[people[0]], b.used[people[1]] = true, true
	// Swap the second and third letters: "Page" -> "Pgae", "Holm" -> "Hlom".
	s := []byte(pair[0])
	s[1], s[2] = s[2], s[1]
	l := b.newList(people, [2]string{first + " " + string(s), people[1]})
	l.items("plain")
	return l.finish("typo-list")
}
