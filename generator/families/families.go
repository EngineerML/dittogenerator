package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation. Only the families that sat near the
// threshold on the frontier panel (Hades v26, horizon v9 and drp23v15, two
// repeats, 2026-11-05) stay in it, and alias-list, whose look-alike
// distractors decide between half credit and none, gets two slots.
//
// Retired because every frontier agent solved them, and kept for experiments:
// AsOf (30/30), Handover (all), LiveStatus (all) and RevertOriginal (34/36).
var Default = []Family{
	AliasList, RevertPrevious, StruckEntry, MisplacedCue, MembershipFold, EventOrder,
	AliasList, ForwardedOld, NicknameEmail, InitialsContact, OwnerFirm,
	Decided, TypoContact,
}
