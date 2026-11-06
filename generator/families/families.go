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
var Default = []Family{
	RoleList, AliasList, RoleList, AliasList, TypoList, RoleList,
	AliasList, RoleList, AliasList, RoleList, AliasList,
}
