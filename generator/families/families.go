package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation: half-credit lists only (lists.go), with the
// links that kept full credit rarest on the frontier panel (Hades v26,
// horizon v9 and drp23v15, 2026-11-05) in the most slots.
//
// With the live champion (dan0l v42) solving almost every case, a case needs
// only one more owner at 0.5, and every full-credit answer raises its mean.
// The pass/fail families were about 60% of Hades v26's full answers on a v6
// pack; without them the estimated live score rises from about 0.31 to 0.33
// (tools/panel_estimate.py, five-member panel). Since alias-list and
// typo-list stopped naming the shop kind in the second record (2026-11-05),
// the Hades line passes them about half the time with almost no full credit.
// Estimated credit per case: alias-list 0.61, role-list 0.56, typo-list 0.50
// (one slot, for variety); about 0.35 for the pack.
//
// Kept for experiments, out of the rotation: the pass/fail families that
// split the panel best (TypoContact, MembershipFold, NicknameEmail,
// MisplacedCue), InitialsList (the Hades line follows initials), and the
// families every frontier agent solved or that earned little (AsOf, Handover,
// LiveStatus, RevertOriginal, RevertPrevious, StruckEntry, EventOrder,
// ForwardedOld, InitialsContact, OwnerFirm, Decided).
var Default = []Family{
	RoleList, AliasList, RoleList, AliasList, TypoList, RoleList,
	AliasList, RoleList, AliasList, RoleList, AliasList,
}
