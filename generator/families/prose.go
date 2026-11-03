package families

// Text the renderers wrap records in. None of it holds a name, number, date
// or anything a question asks about, so it can sit next to any record.

var acks = []string{
	"Got it.", "Thanks, I'll keep that in mind.", "Okay, noted.", "Understood.", "Thanks for the update.",
	"Makes sense.", "All right, I'll remember that.", "Good to know.",
}

var subjects = []string{"Quick update", "Plans", "Re: logistics", "Following up", "FYI", "Schedule"}

var openers = []string{
	"Notes from today's {topic} sync.",
	"Quick write-up on the {topic}.",
	"Here is where things are with the {topic}.",
	"Catching up on the {topic} after the call.",
}

var closers = []string{
	"That's all for the {topic} today.",
	"More on the {topic} next week.",
	"Ping me with questions about the {topic}.",
	"I'll keep this thread going for the {topic}.",
}

var paragraphs = []string{
	"We spent most of the call on logistics and who needs what by when. The action items are in the shared folder.",
	"Everyone agreed the budget sheet needs a cleaner layout before anyone outside the group sees it.",
	"There was a long side discussion about parking, which went in circles, so I've left it out of these notes.",
	"Reminder that the shared calendar is the source of truth for meeting times; please keep it current.",
	"The catering feedback from last time was mostly positive, though the coffee ran out early again.",
	"Someone asked whether we should record the sessions. We'll revisit that once the agenda settles.",
	"The printer in the meeting room is still acting up, so bring copies if you need them.",
	"We also touched on travel reimbursements; the form is in the usual place and takes a few minutes.",
	"A couple of people will join remotely, so the room needs the conference speaker set up beforehand.",
	"The slides from the last review are attached for anyone who missed it.",
	"Weather might be a factor for anything outdoors, so a backup plan is worth sketching out.",
	"Thanks to everyone who sent comments on the draft; most of them are folded in now.",
	"The group chat has been busy, so these notes pull the useful parts into one place.",
	"There's still an open question about name badges, which can wait until the guest count is firm.",
}
