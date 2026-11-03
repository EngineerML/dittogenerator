package families

import (
	"fmt"
	"strings"
	"unicode"
)

// Grammar is a small Tracery-style grammar: each key maps to alternatives,
// and "#key#" inside an alternative expands recursively. Values are filled in
// after expansion (see fill), so a phrasing can never alter a value.
type Grammar map[string][]string

// Expand returns one expansion of s.
func (b *Builder) Expand(g Grammar, s string) string {
	for depth := 0; depth < 8 && strings.Contains(s, "#"); depth++ {
		var out strings.Builder
		parts := strings.Split(s, "#")
		for i, part := range parts {
			alts, ok := g[part]
			if i%2 == 1 && ok {
				out.WriteString(alts[b.Rng.Intn(len(alts))])
				continue
			}
			if i%2 == 1 {
				out.WriteString("#")
			}
			out.WriteString(part)
		}
		s = out.String()
	}
	return s
}

// Phrase expands one of the alternatives for key in g and fills its values.
func (b *Builder) Phrase(g Grammar, key string, kv ...string) string {
	return capFirst(fill(b.Expand(g, "#"+key+"#"), kv...))
}

// fill replaces {name} placeholders with values.
func fill(s string, kv ...string) string {
	pairs := make([]string, 0, len(kv))
	for i := 0; i+1 < len(kv); i += 2 {
		pairs = append(pairs, "{"+kv[i]+"}", kv[i+1])
	}
	return strings.NewReplacer(pairs...).Replace(s)
}

func capFirst(s string) string {
	for i, r := range s {
		return string(unicode.ToUpper(r)) + s[i+len(string(r)):]
	}
	return s
}

// Subject is the thing a case asks about, named the way its path needs.
type Subject struct {
	Path Path
	// Name is how records refer to it; Ask is how the question does.
	Name, Ask string
	// Topic is the phrase the topic path's notes share with the question.
	Topic string
	// File is the case-file identifier on the program path.
	File string
}

// Subject names thing (an event, a list, a contract) for the given path:
// someone's thing on the records path, a thing in a case file on the program
// path, and a project's thing on the topic path.
func (b *Builder) Subject(path Path, thing string) Subject {
	switch path {
	case PathProgram:
		id := b.CaseFile()
		return Subject{Path: path, Name: "the " + thing, Ask: "the " + thing + " in case file " + id, File: id}
	case PathTopic:
		topic := b.Pick(projects, 1)[0] + " " + thing
		return Subject{Path: path, Name: "the " + topic, Ask: "the " + topic, Topic: topic}
	default:
		name := b.Person() + "'s " + thing
		return Subject{Path: path, Name: name, Ask: name}
	}
}

// Record seeds one record about s in a wave: a long multi-paragraph note on
// the topic path, a case-file note on the program path, and a short chat,
// email, text or table row otherwise. Words of protect (values, names) are
// never noised.
func (b *Builder) Record(s Subject, wave int, line string, protect ...string) string {
	// A fresh slice: appending to the caller's would overwrite whatever its
	// backing array holds past the values passed in.
	protect = append(append([]string(nil), protect...), s.Name, s.Topic)
	switch s.Path {
	case PathTopic:
		return b.LongNote(wave, s.Topic, line, protect...)
	case PathProgram:
		return b.Note(wave, "Case file "+s.File+": "+line, protect...)
	default:
		return b.Note(wave, line, protect...)
	}
}

var contextLines = []string{
	"Kickoff notes for {s} are in the shared folder.",
	"Budget questions on {s} go to the usual channel.",
	"{s} has its own thread now, so updates stay together.",
	"Reminder to keep {s} on the weekly agenda.",
}

// Context seeds n records about s that carry nothing a question asks. The
// agents take the program and topic paths only when at least three records
// share the identifier or phrase, so a family with fewer records of its own
// tops them up here.
func (b *Builder) Context(s Subject, n int) {
	for ; n > 0; n-- {
		b.Record(s, b.Rng.Intn(Waves), capFirst(fill(b.Any(contextLines), "s", s.Name)))
	}
}

// Note seeds a short record in one of four renderings.
func (b *Builder) Note(wave int, line string, protect ...string) string {
	line = b.noise(line, protect)
	ack := b.Any(acks)
	switch r := b.Rng.Intn(10); {
	case r < 5:
		return b.Say(wave, line, ack)
	case r < 7:
		from := b.Any(firstNames)
		return b.Say(wave, fmt.Sprintf("Forwarding an email from %s.\n\nSubject: %s\n\n%s\n\nThanks,\n%s", from, b.Any(subjects), line, from), ack)
	case r < 9:
		return b.Say(wave, fmt.Sprintf("Forwarding a text from %s: \"%s\"", b.Any(firstNames), line), ack)
	default:
		return b.Say(wave, "Pasted from the shared planning sheet: "+line, ack)
	}
}

// LongNote seeds a multi-paragraph note that names topic in its opening and
// closing paragraphs and carries line in the middle, among unrelated
// paragraphs.
func (b *Builder) LongNote(wave int, topic, line string, protect ...string) string {
	paras := []string{fill(b.Any(openers), "topic", topic)}
	body := []string{b.noise(line, protect)}
	for n := 1 + b.Rng.Intn(2); n > 0; n-- {
		body = append(body, b.Any(paragraphs))
	}
	b.Rng.Shuffle(len(body), func(i, j int) { body[i], body[j] = body[j], body[i] })
	paras = append(paras, body...)
	paras = append(paras, fill(b.Any(closers), "topic", topic))
	return b.Say(wave, strings.Join(paras, "\n\n"), b.Any(acks))
}

// noise swaps two inner letters of at most one plain lowercase word, the
// kind of slip real notes carry. Words of protected values, names and
// identifiers are never touched, so every value stays exactly as graded.
func (b *Builder) noise(line string, protect []string) string {
	if b.Rng.Intn(4) != 0 {
		return line
	}
	keep := map[string]bool{}
	for _, p := range protect {
		for _, w := range strings.Fields(strings.ToLower(p)) {
			keep[strings.Trim(w, ".,;:'\"()")] = true
		}
	}
	words := strings.Split(line, " ")
	var candidates []int
	for i, w := range words {
		if len(w) >= 6 && !keep[w] && isLowerWord(w) {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return line
	}
	i := candidates[b.Rng.Intn(len(candidates))]
	w := []byte(words[i])
	j := 1 + b.Rng.Intn(len(w)-3)
	w[j], w[j+1] = w[j+1], w[j]
	words[i] = string(w)
	return strings.Join(words, " ")
}

func isLowerWord(w string) bool {
	for _, r := range w {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

// day returns the calendar day offset days after the pack's first day,
// written in one of two plain forms.
func (b *Builder) day(offset int) string {
	d := b.epoch.AddDate(0, 0, offset)
	if b.Rng.Intn(2) == 0 {
		return fmt.Sprintf("%d %s", d.Day(), d.Month())
	}
	return fmt.Sprintf("%s %d", d.Month(), d.Day())
}

// dayFixed returns a day in "14 March" form, for questions.
func (b *Builder) dayFixed(offset int) string {
	d := b.epoch.AddDate(0, 0, offset)
	return fmt.Sprintf("%d %s", d.Day(), d.Month())
}

// list joins values as prose: "A, B and C".
func list(values []string) string {
	if len(values) < 2 {
		return strings.Join(values, "")
	}
	return strings.Join(values[:len(values)-1], ", ") + " and " + values[len(values)-1]
}

// without returns values minus drop.
func without(values []string, drop ...string) []string {
	var out []string
	for _, v := range values {
		keep := true
		for _, d := range drop {
			if v == d {
				keep = false
			}
		}
		if keep {
			out = append(out, v)
		}
	}
	return out
}

// email builds a work address for a person at a domain.
func email(person, domain string) string {
	parts := strings.Fields(strings.ToLower(person))
	return strings.Join(parts, ".") + "@" + domain
}
