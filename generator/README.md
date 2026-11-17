# horizon generator

A deterministic Go program for the DittoBench generator competition. It reads one generator request
on stdin and writes one challenge pack on stdout: memory waves for one user and memory cases whose
answers only a careful reader of those memories gets right.

`main.go` and the builder's shape come from the public starter kit (MIT, see LICENSE.starter-kit);
the families are ours.

## Layout

| File | What |
|---|---|
| `main.go` | entry point; declares the newest grading policy the request accepts that the generator knows, never below the round's own baseline |
| `families/builder.go` | eight waves, opaque pair and session IDs, reserved values and people, timestamps |
| `families/text.go` | a small grammar, renderers (chat, forwarded email, forwarded text, sheet row, long note) and letter-swap noise that never touches a value |
| `families/multi.go` | the multi-link list family in the rotation (v8) |
| `families/edges.go` | link wording, name shapes, four-digit numbers, person chatter |
| `families/lists.go` | earlier list families and the shop kinds |
| `families/chains.go`, `joins.go`, `splitters.go` | single-answer families, kept for experiments |
| `families/pools.go`, `prose.go` | values and the text around them |

## Families

Every case asks for the two places one person keeps going to. One is told under the person's full
name (with she or he), the other only through a link the reader has to follow, and every link brings
a look-alike place that is not the person's. A list answer is graded by the share of items named,
and naming a look-alike place or declining zeroes it, so the full-name place earns half credit,
following the link earns full credit, and taking a look-alike earns nothing.

| Link | The second place is told under | Look-alike |
|---|---|---|
| `session` | the next turn of a conversation that named the person, with another conversation in between | another person's conversation under a different lead |
| `number` | receipts under a four-digit membership number another record gives as the person's | a receipt under the same digits in another order |
| `title` | Ms. or Mr. and the surname, whichever fits the full-name record's pronoun | the surname under the other title |
| `shape` | a name outside the two-capitalised-words shape (double first names, particles, Mac/O' forms) | a person the two-word reading confuses with them |
| `found` | the place the person moved to after giving one up | the place given up |

Most cases add one more look-alike: a relative under the other title (`title`) or a place the person
used and later gave up (`givenup`). Hard forms say it loosely (first name only, no stated tie);
soft forms name the person in full and say the tie ("Dana Reyes's brother, Mr. Reyes"). The
rotation mixes hard and soft so the cases split strong and weaker readers alike. Every link and
every giving-up is stated plainly; filler chatter about the person surrounds the records.

## Checks

```bash
go vet ./... && go test ./...      # every size 8..80 validates on several seeds; pools hold the rotation
```

The competition's `challenge rehearse` builds and runs the module the way the scorer does, with
validation and a determinism re-run.
