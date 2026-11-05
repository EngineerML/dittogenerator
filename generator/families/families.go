package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation: ten of its thirteen slots are half-credit
// list families (lists.go), the other three the pass/fail families that split
// the frontier panel most (Hades v26, horizon v9 and drp23v15, 2026-11-05),
// kept so one reading fix cannot undercut the whole pack.
//
// With the live champion (dan0l v42) solving almost every case, a list case
// needs only one more owner at 0.5 and keeps a low mean. Estimated credit per
// case (tools/panel_estimate.py, five-member panel): role-list 0.56,
// alias-list 0.54, typo-list 0.51, typo-contact 0.45, membership-fold 0.42,
// nickname-email 0.34; about 0.31 for the pack.
//
// Out of the rotation and kept for experiments: InitialsList (the Hades line
// follows initials, 0.38), MisplacedCue (0.28-0.38), AsOf, Handover,
// LiveStatus and RevertOriginal (solved by every frontier agent), and
// RevertPrevious, StruckEntry, EventOrder, ForwardedOld, InitialsContact,
// OwnerFirm and Decided (0.13-0.37).
var Default = []Family{
	RoleList, AliasList, TypoList, RoleList, AliasList, TypoContact, RoleList,
	AliasList, MembershipFold, RoleList, AliasList, TypoList, NicknameEmail,
}
