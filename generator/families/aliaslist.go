package families

import (
	"fmt"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// AliasList is a two-item list about one person whose second item is told
// only under a nickname, next to a namesake with a nickname and items of
// their own. It is the shape behind the generator competition's half-credit columns
// (docs/v14-generator-credits.md, SN-352): a list is graded by the share of
// items named, and any distractor zeroes it. An agent that reads the full
// name exactly names one item (0.5, which counts as solved); one that matches
// the first name, or joins the wrong nickname, names a namesake's item (0);
// one that follows the identity record names both (1). On the live miner
// column built this way, horizon v8 scored 0.5 on 23-26 of 35 cases and the
// panel means sat near 0.2-0.5.
//
// Both nickname bindings are stated once in plain words, so the case is
// solvable at full credit.
var aliasLines = Grammar{
	"bind": {
		"{p} goes by {n} with pretty much everyone.",
		"Everyone calls {p} {n}, has done for years.",
		"Quick note so I don't mix them up: {p} is the one we call {n}.",
	},
	"full": {
		"{p} loves {v} and goes there most weekends.",
		"{p} keeps telling me {v} is the best of the {c} around here.",
		"Ran into {p} at {v} again; it's the favourite of all the {c} for sure.",
	},
	"nick": {
		"{n} swears by {v}, says none of the other {c} come close.",
		"{n} dragged me to {v} today and would not stop praising it.",
		"Apparently {v} is where {n} goes every single week.",
	},
	"ask": {
		"Which two {c} does {p} love?",
		// Questions only: the imperative "Name the two..." sent every frontier
		// agent to its action lane.
		"Which two {c} is {p} really fond of?",
		"What are the two {c} {p} keeps going back to?",
	},
}

// aliasKinds are the shop kinds an item is drawn from: the plural used in
// the text and the noun that ends each shop's name.
var aliasKinds = [][2]string{
	{"bakeries", "Bakehouse"}, {"bookshops", "Books"}, {"coffee roasters", "Roasters"},
	{"delis", "Deli"}, {"florists", "Flowers"}, {"bike shops", "Cycles"},
}

// AliasList seeds the two bindings and four items (two per person, one under
// the full name and one under the nickname), each in its own wave, and asks
// for the first person's two.
func AliasList(b *Builder, i int) challenge.Case {
	first := b.Any(firstNames)
	var people []string
	for len(people) < 2 {
		if name := first + " " + b.Any(surnames); !b.used[name] {
			b.used[name] = true
			people = append(people, name)
		}
	}
	nick := b.Pick(nicknames, 2)
	kind := aliasKinds[b.Rng.Intn(len(aliasKinds))]
	shops := b.Pick(shopNames, 4)
	for k := range shops {
		shops[k] += " " + kind[1]
	}
	w := b.Rng.Perm(Waves)[:6]

	say := func(wave int, key string, kv ...string) string {
		line := capFirst(b.Phrase(aliasLines, key, append(kv, "c", kind[0])...))
		return b.Note(wave, line, append(append([]string(nil), people...), append(nick, shops...)...)...)
	}
	ids := []string{
		say(w[0], "bind", "p", people[0], "n", nick[0]),
		say(w[1], "full", "p", people[0], "v", shops[0]),
		say(w[2], "nick", "n", nick[0], "v", shops[1]),
	}
	say(w[3], "bind", "p", people[1], "n", nick[1])
	say(w[4], "full", "p", people[1], "v", shops[2])
	say(w[5], "nick", "n", nick[1], "v", shops[3])

	gold := shops[:2]
	return challenge.Case{
		MemoryCase: protocol.MemoryCase{
			// The rotation holds two alias-list slots, so the round number alone
			// would repeat; the case's place in the pack does not.
			ID:                fmt.Sprintf("alias-list-%d", b.Len()),
			QuestionType:      "alias-list",
			Question:          b.Phrase(aliasLines, "ask", "c", kind[0], "p", people[0]),
			ExpectedAnswer:    strings.Join(gold, ", "),
			AnswerKind:        protocol.AnswerList,
			AnswerItems:       gold,
			DistractorAnswers: shops[2:],
		},
		EvidencePairIDs: ids,
	}
}
