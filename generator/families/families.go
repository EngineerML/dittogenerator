package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation: nine of its thirteen slots are the
// half-credit list families (lists.go), four links to the second item so no
// single reading habit beats them all; the other four are the pass/fail
// families that split the frontier panel most (Hades v26, horizon v9 and
// drp23v15, 2026-11-05).
//
// With the live champion (dan0l v42) solving almost every case, a list case
// needs only one more owner at 0.5 and keeps a low mean, so it earns about
// twice the credit of a pass/fail case (tools/panel_estimate.py).
//
// Out of the rotation and kept for experiments: AsOf, Handover, LiveStatus and
// RevertOriginal (solved by every frontier agent), and RevertPrevious,
// StruckEntry, EventOrder, ForwardedOld, InitialsContact, OwnerFirm and
// Decided (0.13-0.37 estimated credit per case, against 0.44-0.51 for the
// families kept).
var Default = []Family{
	AliasList, InitialsList, RoleList, TypoList, TypoContact,
	AliasList, InitialsList, MembershipFold, RoleList, TypoList,
	NicknameEmail, AliasList, MisplacedCue,
}
