package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation: chain families first, then indirection,
// then the wording splitters, so any pack size keeps the mix close to
// 40/35/25.
var Default = []Family{
	RevertPrevious, RevertOriginal, StruckEntry, MisplacedCue, MembershipFold, EventOrder,
	AsOf, ForwardedOld, NicknameEmail, InitialsContact, OwnerFirm,
	Decided, LiveStatus, TypoContact, Handover,
}
