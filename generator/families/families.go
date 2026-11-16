package families

import "github.com/ditto-assistant/dittobench-datagen/challenge"

// Family emits one case into the builder. i is the family's round in the
// pack; it picks the variant and the path.
type Family func(b *Builder, i int) challenge.Case

// Default is the family rotation: list families only (lists.go). The links
// readers follow least reliably, a role and a nickname, take most of the
// slots; the misspelled name takes one, for variety.
//
// Kept for experiments, out of the rotation: InitialsList, and the
// single-answer families in chains.go, joins.go and splitters.go.
//
// v8 (2026-11-16): edge families (edges.go) take most slots, each sitting just
// past a limit many current agents share; AliasList (nickname bound inside a
// long note) and RoleNoteList keep the link families. See notes/
// 2026-11-16_edge_tracks.md.
//
// v8.1: the two-item edge lists came out trivial on a released-agent panel (25 of
// 35 cases solved by all four at 0.5), so the rotation is MultiList only:
// three-item lists over two different edge links (multi.go).
var Default = []Family{
	MultiList, MultiList, MultiList, MultiList, MultiList, MultiList,
	MultiList, MultiList, MultiList, MultiList, MultiList,
}
