# Generator competition work

Our entry for SN118's generator competition: a deterministic Go program that writes memory
challenge packs the top agents split on. Rules, scoring and field evidence:
notes/2026-10-27_generator_competition.md.

Branch `generator`, cut from harness/dan0l-combined so the agent-side failure notes and tools
are at hand.

## Layout

`main.go` and the builder's shape come from the public starter kit
(`repos/ditto-subnet/miners/dittobench-generator-starter-kit`, MIT, see LICENSE.starter-kit). The
rest is ours (version 3):

| File | What |
|---|---|
| `families/builder.go` | eight waves, opaque pair and session IDs, reserved values and people, case-file IDs, timestamps |
| `families/text.go` | a small grammar (`#key#` alternatives, values filled in after), subjects per path, renderers (chat, forwarded email, forwarded text, sheet row, long multi-paragraph note), letter-swap noise that never touches a value |
| `families/chains.go`, `joins.go`, `splitters.go` | the fifteen families |
| `families/pools.go`, `prose.go` | values and the text around them |

The packs declare the newest grading policy the request accepts that this generator knows (13 to
15): 15 in a bench-15 or bench-16 round. A bench-16 round accepts 15 and 16, and policy 16 grades
exactly like 15, so declaring 16 would change nothing.

### Paths

The top agents route a memory question on surface features before they read anything: a long
identifier carried by at least three records sends it to their case-file (program) path, a phrase
shared with at least three multi-paragraph notes to their topic path, a full name to the records
lane. Each path reads chains with its own rules. Most families rotate one case per round over the
three paths (`records`, `program`, `topic`), and the path is the last part of the case's family
name in every summary.

### Families

Built from the agent research in /home/bittensor/ditto (v15_atlas, notes 2026-11-01 to 2026-11-03):
the shapes there that gpt-oss-20b agents get right 20 to 80 percent of the time, and the places the
DrPresident and dan0l lines route the same question differently.

| Group | Family | What happens | Gold |
|---|---|---|---|
| chain | `revert-previous` | three or four moves, then "reverts to the previous entry" (mid) or a plainer undo | the city before the last move |
| chain | `revert-original` | two or three moves, then "back to the city we had at the very start" | the first city |
| chain | `struck-entry` | the location worded differently in each entry, then the latest entry "logged in error" | the city before it |
| chain | `misplaced-cue` | a draft "update" that sign-off turns down, or a plain move followed by a declined proposal | the city that stands |
| chain | `membership-fold` | four people, then a removal, an addition, a swap (and a re-add), asked "after the change" | the final list |
| chain | `event-order` | five milestones, two dated only relative to another, filed in random waves | the one right before a named milestone |
| join | `as-of` | three dated locations | the one that held on a day inside the middle interval |
| join | `forwarded` | a dated move, then an older message forwarded last | the moved-to city |
| join | `nickname-email` | an email changed in a note that uses only the person's nickname | the new email |
| join | `initials-contact` | a thread handed to "C. Nilsen", with a directory holding a namesake | the new holder's email |
| join | `owner-firm` | a contract moved to a firm whose coordinator is named in the same note, asked "who holds it" | the firm |
| splitter | `decided` | two cities weighed, a third suggested, one agreed; asked "remind me what we decided" | the agreed city |
| splitter | `live-status` | three multi-word states; asked for "the up-to-date status" | the latest state |
| splitter | `typo-contact` | a handover in long notes, emails on short cards; asked "whats the emial" (or cleanly) | the new runner's email |
| splitter | `handover` | a thread handed over with "took over", "run point", "has the ball" or "owns", and a third person cc'd | the new holder |

Every wrong value a case's memory holds is listed as a distractor, so a reply that names it scores
0 and a memory dump fails.

## Checks

```bash
cd generator && go vet ./... && go test ./...      # every size 8..80 validates on several seeds
cd repos/ditto-subnet/research/dittobench-datagen
go run ./cmd/challenge rehearse -dir ../../../../generator \
  -hotkey 5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY -seeds 3 -cases 80 -round-bench 15
```

Rounds have opened at bench 15 since 2026-11-03 (packs graded at 14 or 15), and the bench 16
(dueling) rollout opened on 2026-11-04 (15 or 16). Rehearse with `-round-bench 15`, which grades
like a bench-16 round; the packs then declare 15. In bench 16 the top generators' columns are also
35 cases each of every agent's composite, so their source is published after an embargo.
`challenge-scorer run -dir`
sends the original v13 request, so the panel tool rehearses at the round's bench with `-out-dir`
and runs each pack with `challenge-scorer run -pack`.

`tools/generator_panel.py` runs the packs against a local panel of agent images with the
official challenge scorer and prints, per family and path, each agent's solved/n and mean and the
panel's discriminating/trivial/unsolved split (paid: OpenRouter, about $0.003 per case per agent).
With `--mock` it is free: a fake model and fake embeddings stand in, and the traced images record
per case the lane, the reading path and whether the evidence records reached the model.

```bash
python3 tools/generator_panel.py --mock --agent drp=drp23v15 --agent dan0l=v30-released --seeds 1
```

## Results against our own agent (horizon v2 image, one seed, 80 cases)

| Version | Family | Solved |
|---|---|---|
| v1 | elimination / swapped / picked | 20/20, 20/20, 19/20 (trivial: replaced) |
| v1 | revert ("back to the very beginning" / "undo the last change") | 13/20 (misses are abstentions) |
| v2 | revert / position | 16/20, 16/20 |
| v2 | same-as / undo-revert | 20/20, 20/20 |

About $0.25 per 80-case run. One agent cannot show what is discriminating (that needs two
owners solving and at least one agent failing); the next step is a local panel of different
agents.
