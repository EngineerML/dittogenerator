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
| `families/lists.go` | the list families in the rotation |
| `families/chains.go`, `joins.go`, `splitters.go` | single-answer families, kept for experiments |
| `families/pools.go`, `prose.go` | values and the text around them |

## Families

Every case asks for the two shops one person likes. One is told under the person's full name, the
other only through a link the reader has to follow, and a look-alike person nearby has shops of
their own. A list answer is graded by the share of items named, and naming a look-alike's shop or
declining zeroes it, so naming the full-name shop earns half credit, following the link earns full
credit, and confusing the two people earns nothing. Every link is stated plainly or unambiguous.

| Family | Share | The second shop is told under | Look-alike |
|---|---|---|---|
| `role-list` | 5 of 11 | a role ("the Brevant lead"), named in another record | a namesake leading another team |
| `alias-list` | 5 of 11 | a nickname, bound once inside a note about something else | a namesake (same first name) with a nickname |
| `typo-list` | 1 of 11 | the name with two surname letters swapped ("Edith Pgae") | a real surname one letter away ("Edith Pace") |

The second record of `alias-list` and `typo-list` never names the shop kind the question asks about,
so the question's words alone do not lead to it.

## Checks

```bash
go vet ./... && go test ./...      # every size 8..80 validates on several seeds; pools hold the rotation
```

The competition's `challenge rehearse` builds and runs the module the way the scorer does, with
validation and a determinism re-run.
