package families

// senderNames sign forwarded emails and texts. No case person can carry one of
// them as a first name (they are removed from firstNames' reach below), so a
// note "from Beatrix" never sits next to a question about Beatrix Popescu.
var senderNames = func() []string {
	pool := []string{
		"Wilhelmina", "Thaddeus", "Ottilie", "Barnaby", "Henrietta", "Ignatius", "Clementine",
		"Leopold", "Philippa", "Rupert", "Seraphina", "Cornelius", "Winifred", "Augustin",
		"Marguerite", "Ambrose", "Rosalind", "Percival", "Georgiana", "Lysander",
	}
	taken := map[string]bool{}
	for _, f := range firstNames {
		taken[f] = true
	}
	var out []string
	for _, s := range pool {
		if !taken[s] {
			out = append(out, s)
		}
	}
	if len(out) < 8 {
		panic("sender pool overlaps the person first names")
	}
	return out
}()
