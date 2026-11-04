package families

import (
	"fmt"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

// Restated is the benchmark's parser-divergence shape carried across waves,
// the one memory family every agent line still misses about two times in
// three (local benchmark runs of DrPresident v23 to v23-v15, dan0l v30 and
// horizon v8, 2026-10-31 to 11-03).
//
// Each record holds its own trap, a value it says is wrong (negated, only
// reported by someone else, or hypothetical) next to the right one, and the
// same fact is told again in a later wave with a different value and no
// word saying it changed. The same question, word for word, is asked after
// the middle telling and again after the last wave; each time the newest
// telling seen so far is the answer. The agents answer every copy with one
// value, so at most one is right: they rank the identical-looking records by
// relevance, not by when they were said, and may read back their own
// earlier reply.
// A word follows every number: the grader reads a "." or "," right after
// digits as part of the number, so "it is 1185." would hold neither the gold
// nor the decoy.
var restatedLines = Grammar{
	"negated": {
		"One thing to keep straight: {s} is not {a} at all. It is {v} and nothing else.",
		"Correcting the record on {s}: it isn't {a} as I wrote somewhere, it's {v} for sure.",
		"For the avoidance of doubt, {s} is {v} and not {a} whatever the old sticky note says.",
	},
	"reported": {
		"{p} insisted {s} was {a} on the phone, but my own note says {v} and my note is the one to trust.",
		"{p} told me {s} is {a} over lunch. I checked myself and it is {v} on the paperwork.",
		"{p} keeps saying {s} is {a} but it is actually {v} according to the building office.",
	},
	"hypothetical": {
		"If the replacement lock had gone in, {s} would have been {a} instead, but that never happened, so it is {v} on the key tag.",
		"{s} would be {a} if the landlord had signed off on it. He didn't, so it is {v} as written on the card.",
		"Had the reset gone through, {s} would read {a} on the panel; it didn't, and it is {v} on the panel.",
	},
	"ask": {
		"What is {a}?",
		"Remind me, what is {a}?",
		"Which number is {a}?",
	},
}

var codeThings = []string{
	"gate code", "locker number", "storage unit code", "alarm code", "mailbox code",
	"garage door code", "bike lock code", "key box code", "padlock code", "door code",
}

// Restated tells a code three times in three waves and asks for it after the
// middle telling and after the last wave, plus after the first telling every
// other round, as long as the pack has room.
func Restated(b *Builder, i int) challenge.Case {
	path := pathFor(i, 2)
	s := b.Subject(path, b.Any(codeThings))
	w := b.Spread(3)
	forms := []string{"negated", "reported", "hypothetical"}
	b.Rng.Shuffle(len(forms), func(x, y int) { forms[x], forms[y] = forms[y], forms[x] })

	right, wrong, ids := make([]string, 3), make([]string, 3), make([]string, 3)
	for k := range w {
		right[k], wrong[k] = b.Code(), b.Code()
		p := b.Any(firstNames)
		line := capFirst(b.Phrase(restatedLines, forms[k], "s", s.Name, "a", wrong[k], "v", right[k], "p", p))
		ids[k] = b.Record(s, w[k], line, right[k], wrong[k], p)
		b.Filler()
	}
	question := b.Phrase(restatedLines, "ask", "a", s.Ask)

	// ask is the case asked right after wave `after`, when tellings 0..last
	// are visible: the newest is the gold, every other visible value a decoy.
	ask := func(stage string, last, after int) challenge.Case {
		decoys := append(append([]string(nil), wrong[:last+1]...), right[:last]...)
		return challenge.Case{
			MemoryCase: protocol.MemoryCase{
				ID:                fmt.Sprintf("restated-%s-%s-%d", stage, path, i),
				QuestionType:      "restated-" + stage,
				Question:          question,
				ExpectedAnswer:    right[last],
				AnswerKind:        protocol.AnswerNumber,
				DistractorAnswers: decoys,
			},
			EvidencePairIDs: []string{ids[last]},
		}
	}
	if i%2 == 0 && b.Spare() >= 2 {
		b.AddAt(ask("first", 0, w[0]), w[0])
	}
	if b.Spare() >= 1 {
		b.AddAt(ask("mid", 1, w[1]), w[1])
	}
	return ask("last", 2, Waves-1)
}
