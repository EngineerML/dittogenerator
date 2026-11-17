package families

// v9 links and traps (2026-11-17). Released generators had shown the title,
// given-up, found-instead and session-relay mechanics to the agents, and their
// columns lost most of their score within a day or two, so v9 asks for a
// different reading.
//
// linkCond: one record says which of two places the person uses depends on a
// condition ("on early shifts it's X; on late shifts it's Y"); another record,
// which never names the place kind, says which condition holds now. Either
// branch may be the one that holds. The place under the other branch is the
// distractor. The question always asks about now ("these days", "at the
// moment"): asked without a time ("keeps going back to", "love"), both
// branches' places are fair answers, and v9's first panel run listed both.
// Soft: one record repeats the branch's own words. Hard: an earlier record
// says the other branch held, and a later one says in other words, still
// unambiguously, that this branch holds now.
//
// trapIntent: a place the person considered or tried once and did not take up.
// Soft: turned down or decided against. Hard: tried for a week, did not go back.

// branch is one side of a condition: its phrase in the conditional record and
// the records saying it holds (plainly, and in other words).
type branch struct {
	phrase string
	plain  []string
	other  []string
}

// condition is one either-or circumstance with the conditional record's topic.
type condition struct {
	topic string
	sides [2]branch
}

var conditions = []condition{
	{"the shift pattern", [2]branch{
		{"on early shifts",
			[]string{"{p} is on early shifts now, until further notice.", "{p} has been moved onto early shifts for the rest of the year."},
			[]string{"{p}'s rota now starts at six in the morning and finishes at two, every day.", "{p} is out of the door before half five every morning for work these days."}},
		{"on late shifts",
			[]string{"{p} is on late shifts now, until further notice.", "{p} has been moved onto late shifts for the rest of the year."},
			[]string{"{p}'s new rota starts at two in the afternoon and runs until ten at night, every day.", "{p} doesn't start work before two in the afternoon any more and finishes around ten."}},
	}},
	{"where the work is", [2]branch{
		{"on office days",
			[]string{"{p} is back in the office every day now.", "{p} works from the office full time again."},
			[]string{"{p}'s firm ended remote working, so it's the commute into town five days a week now.", "{p} has a desk at head office again and is there Monday to Friday."}},
		{"while working from home",
			[]string{"{p} works from home full time now.", "{p} is working from home for good now."},
			[]string{"{p}'s company gave up its office lease, so the whole team is remote for good.", "{p} hasn't been into an office since the firm went fully remote, and won't be."}},
	}},
	{"the car", [2]branch{
		{"when the car is running",
			[]string{"{p}'s car is running again.", "{p} has the car back on the road."},
			[]string{"{p} collected the car from the garage, fixed and back in daily use.", "{p} is driving everywhere again now the gearbox is sorted."}},
		{"while the car is off the road",
			[]string{"{p}'s car is off the road for the next few months.", "{p}'s car will be off the road until spring."},
			[]string{"{p}'s car failed its inspection and the garage can't fix it until spring.", "{p} is waiting months for a gearbox part, so the car just sits in the drive."}},
	}},
	{"the knee", [2]branch{
		{"when the knee is fine",
			[]string{"{p}'s knee is fine again.", "{p}'s knee has healed completely."},
			[]string{"{p} was signed off by the physio and is back to running five miles.", "{p} threw away the knee brace, all healed."}},
		{"while the knee is bad",
			[]string{"{p}'s knee is bad again and will be for a while.", "{p}'s knee is playing up again, probably for months."},
			[]string{"{p} is on crutches for the next six weeks after the knee operation.", "{p} was told to keep weight off the injured leg until the new year."}},
	}},
	{"where home is", [2]branch{
		{"while living in the flat in town",
			[]string{"{p} is living in the town flat now.", "{p} has moved into the flat in town."},
			[]string{"{p} let the cottage go and has the town flat as home now.", "{p}'s post goes to the flat in town these days; the cottage is let to someone else."}},
		{"since moving out to the cottage",
			[]string{"{p} has moved out to the cottage for good.", "{p} lives at the cottage now, the move is done."},
			[]string{"{p} handed back the keys to the town flat; the cottage is home now.", "{p}'s town flat has new tenants, and the cottage is where the post goes now."}},
	}},
	{"the job", [2]branch{
		{"while at the bank",
			[]string{"{p} is working at the bank now.", "{p} has taken the job at the bank."},
			[]string{"{p} spends the working week behind a counter in the high-street branch now, after leaving the wards.", "{p} gave up nursing and counts other people's money for a living these days."}},
		{"since starting at the hospital",
			[]string{"{p} has left the bank and started at the hospital.", "{p} has started at the hospital now, the bank job is over."},
			[]string{"{p}'s first month on the hospital wards has just ended, and there is no going back to counting other people's money.", "{p} handed in the security pass at the bank on Friday and starts ward rounds on Monday."}},
	}},
}

var v9Lines = Grammar{
	"cond": {
		"Where {p} goes for {c} depends on {topic}: {a1} it's {x1}, {a2} it's {x2}.",
		"Which of the {c} {p} uses depends on {topic}. {A1}, {x1}; {a2}, {x2}.",
	},
	"askNow": {
		"Which {c} does {p} go to these days?",
		"Which {c} is {p} going to at the moment?",
	},
	"intentSoft": {
		"{p} was offered a free trial at {v} but turned it down.",
		"{p} looked round {v} once and decided against it.",
	},
	"intentHard": {
		"{p} tried {v} for a week last month and didn't go back.",
		"{p} went to {v} twice in the spring and then stopped bothering.",
	},
}

func (m *multi) linkCond(soft bool) {
	item, trap := m.shop(), m.shop()
	c := conditions[m.b.Rng.Intn(len(conditions))]
	now := m.b.Rng.Intn(2) // the branch that holds now
	hold, gone := c.sides[now], c.sides[1-now]
	// which branch the conditional record states first is drawn apart from
	// which branch holds
	x1, x2, a1, a2 := item, trap, hold.phrase, gone.phrase
	if m.b.Rng.Intn(2) == 1 {
		x1, x2, a1, a2 = trap, item, gone.phrase, hold.phrase
	}
	keep := []string{hold.phrase, gone.phrase}
	w1, w2 := m.two()
	ids := []string{m.note(w1, m.b.Phrase(v9Lines, "cond", "p", m.person, "c", m.kind[0], "topic", c.topic,
		"a1", a1, "A1", capFirst(a1), "x1", x1, "a2", a2, "x2", x2), keep...)}
	pick := func(pool []string) string { return fill(pool[m.b.Rng.Intn(len(pool))], "p", m.person) }
	if soft {
		ids = append(ids, m.note(m.wave(), pick(hold.plain), keep...))
	} else {
		early, late := m.two()
		if w2 > late { // the superseded state may sit anywhere before the current one
			late = w2
		}
		if early >= late {
			early = late - 1
			if early < 0 {
				early, late = 0, 1
			}
		}
		ids = append(ids, m.note(early, pick(gone.plain), keep...), m.note(late, pick(hold.other), keep...))
	}
	m.items, m.traps, m.ids = append(m.items, item), append(m.traps, trap), append(m.ids, ids...)
}

func (m *multi) trapIntent(soft bool) {
	trap := m.shop()
	key := "intentHard"
	if soft {
		key = "intentSoft"
	}
	m.note(m.wave(), m.b.Phrase(v9Lines, key, "p", m.person, "v", trap))
	m.traps = append(m.traps, trap)
}
