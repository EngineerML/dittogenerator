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
var Default = []Family{
	NameShapeList, CorrectionList, SessionRelay, NumberList, AliasList, NameShapeList,
	RoleNoteList, CorrectionList, SessionRelay, NameShapeList, NumberList,
}
