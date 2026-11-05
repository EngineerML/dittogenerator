# Generator competition work

Our entry for SN118's generator competition: a deterministic Go program that writes memory
challenge packs the top agents split on. Rules, scoring and field evidence:
notes/2026-10-27_generator_competition.md.

Branch `generator`, cut from harness/dan0l-combined so the agent-side failure notes and tools
are at hand.

## Layout

`main.go` and the builder's shape come from the public starter kit
(`repos/ditto-subnet/miners/dittobench-generator-starter-kit`, MIT, see LICENSE.starter-kit). The
rest is ours (version 6):

| File | What |
|---|---|
| `families/builder.go` | eight waves, opaque pair and session IDs, reserved values and people, case-file IDs, timestamps |
| `families/text.go` | a small grammar (`#key#` alternatives, values filled in after), subjects per path, renderers (chat, forwarded email, forwarded text, sheet row, long multi-paragraph note), letter-swap noise that never touches a value |
| `families/lists.go` | the four half-credit list families |
| `families/chains.go`, `joins.go`, `splitters.go` | the pass/fail families (four in the rotation, the rest kept for experiments) |
| `families/pools.go`, `prose.go` | values and the text around them |

The packs declare the newest grading policy the request accepts that this generator knows (13 to
15), and never less than the round's own baseline: 15 in a bench-15 or bench-16 round, 16 in a
bench-17 round (which accepts 16 and 17). Policies 16 and 17 grade exactly like 15, so the
generator needs no change to stay accepted as rounds move up.

### Paths

The top agents route a memory question on surface features before they read anything: a long
identifier carried by at least three records sends it to their case-file (program) path, a phrase
shared with at least three multi-paragraph notes to their topic path, a full name to the records
lane. Each path reads chains with its own rules. Most families rotate one case per round over the
three paths (`records`, `program`, `topic`), and the path is the last part of the case's family
name in every summary.

### Families

Measured on a frontier panel (Hades v26, horizon v9 built from its upload, drp23v15; official
grader, policy 15) and weighed with `tools/panel_estimate.py`, which models the live champion
(dan0l v42, 0.89 on the top generators' shared cases) that the local panel lacks. With a champion
that solves almost everything, a case needs only one more owner at 0.5, so two-item lists, where
the easy half is worth exactly 0.5, earn about twice the credit of pass/fail cases.

Nine of the thirteen rotation slots are lists (about 70% of a pack). Each asks for two shops a
person likes: one told under the full name, one only through a link, beside a look-alike with
shops of their own. Naming only the first earns 0.5; hedging about the second ("and no record of a
second one") or naming the look-alike's earns 0.

| Family | Slots | The second shop is told under | Look-alike |
|---|---|---|---|
| `alias-list` | 3 | a nickname, bound once ("everyone calls Edith Page Pip") | a namesake (same first name) with a nickname |
| `initials-list` | 2 | initial and surname ("E. Page") | same surname, another initial |
| `role-list` | 2 | a role ("the Brevant lead"), named in another record | a namesake leading another team |
| `typo-list` | 2 | the name with two surname letters swapped ("Edith Pgae") | a real surname one letter away ("Edith Pace") |
| `typo-contact` | 1 | pass/fail: handover in long notes, asked "whats the emial" | |
| `membership-fold` | 1 | pass/fail: removal, addition, swap and re-add, asked "after the change" | |
| `nickname-email` | 1 | pass/fail: an email changed in a note that quotes only the nickname | |
| `misplaced-cue` | 1 | pass/fail: a declined proposal (a rejected draft every third round) | |

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
