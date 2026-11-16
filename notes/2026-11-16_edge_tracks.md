# 2026-11-16 Edge tracks for the generator

From 2026-11-16 the user gave full permission on this folder. Their goal is a generator in the board's top 5, built on edge tracks: general behaviour weaknesses that current agents share.

## Scoring (research/dittobench-datagen/challenge/score.go, live default params)
- Panel knobs: MinPanel 3, MinSolverOwners 2, SolveThreshold 0.5, ChampionBonus 0.5, TargetCases 32.
- Case statuses, among agents not owned by the generator's owner:
  - solved: score ≥ 0.5, so half credit on a 2-item list counts as solved;
  - trivial: solved by everyone;
  - discriminating: solved by ≥ 2 owners and not trivial;
  - unsolved: not solvable and not frontier.
- credit = (1 − mean) × (1.5 if the champion failed).
- pack = Σ credit / (1.5 · max(n, 32)) × (1 − unsolved share)².
- The run score is the mean over seeds. Fingerprints are claimed across seeds, so output must be fresh per seed.
- Emission: the generator track is 3000 bps of the track pool (≈ 24 % of emission). The top 5 hotkeys split it 40/25/15/12/8 % by best latest round score.
- Pins come from the top eligible generators, one per owner. A column owned by the scored agent's coldkey is masked and its 10 % goes back to datagen. So a pinned generator of ours also helps our agent.

## Board (2026-11-16 10:04)
| Generator | Score |
|---|---|
| gan-v1 v2 | 0.248 |
| lh2gen v4 | 0.208 |
| tern-gen v1 | 0.198 |
| 5DWjUM v9 | 0.134 |
| 5D4XN6 v8 | 0.133 (5th) |
| querverweis v3 | 0.120 |
| bauni v6 | 0.039 |

## Column survey (ditto `live/edge_survey.sh`: 128 agents scored on bench 18 since 11-14, per-column means joined with pins)
- **gan-v1 v2** (#1): every agent 0.61–0.71, sd 0.04. Nearly every case is discriminating with credit ≈ 0.35–0.4. Inferred shape: an easy half, a hard half, and a trap that zeroes some agents.
- **tetgol v6:** mean 0.28, p90 0.57 (crown 0.59, most about 0.2), solved ≥ 0.5 by 16/72 agents. Number-keyed chains.
- **lh2gen v4:** mean 0.33; only lh2mem tops 0.6.
- **a1d97bc9 (tessera v2?):** 0.03 to 0.97. Its source was released 11-16 10:21.
- **memgen v7:** mean 0.45, 0.27 to 0.99.
- **querverweis v3:** mean 0.46, 0.34 to 0.85.
- **Near-trivial:** sn118-recall, 5D4XN6 v8 and 5FqQq9 v2 all average ≥ 0.7.

## Shared code limits (crown-v228 C, dan0l v59 D, beast-d41 v19 B; read-only survey)
**Size**
- 60 records per packet (all three).
- A record over the packet budget is skipped whole (all three).
- Packet budget is 24k chars (C, D) and 16k (B).
- Over 30k chars, or two records over 7k, switches the packet to first-line search mode.
- Records over 2,500 chars follow no name edges.
- Link-hops only reach records of 700 chars or less, 2 hops.
- Nickname cards are only read from records of 400 chars or less.

**Time**
- C follows same-session turns for 20 minutes only, at most 4 per record and 6 in all.
- D follows later turns only for "which two..five" asks: the next 3 records, same session, no other session in between.
- B has no session logic.
- Date asks need "Month D, YYYY".

**Names and identifiers**
- `FULL_NAME_RE` is exactly two `[A-Z][a-z]+` words. It misses McDonald, O'Neil, accents, particles, "Mary Ann", "Smith-Jones".
- C's initials pattern needs a dot and a surname of 3 or more letters.
- Quoted nicknames need double quotes and a capital letter.
- `CODE_RE` needs a dash plus letters and digits.
- D's numbers need 5 or more digits and are read from the prompt only; C has no number chain.

**Chains**
- Retrieval follows 3 hops; a value held by more records than its fanout cap is dropped.
- B's reference chain is 4 rounds; D's number chain is 5 hops and 2–12 holders.

**Supersession**
- D's "<X> replaced <Y>" misses passive voice, "swapped", "switched".
- The withdraw words miss scratch / ignore / cancel / disregard.
- C's handle-retired wording misses ditched / gave up.

**Fixed assumptions**
- All three present the two latest states as a "disagreement", even when the later one corrects the earlier.

## Design rule
- Each case is a 2-item list: an easy item (half credit for most), a hard item just past one edge, and a trap that zeroes careless readers. That keeps cases solvable by ≥ 2 owners while the mean stays at about 0.5–0.65.
- Mix several edges per pack: a code edge can be patched within hours of a source release.
- Weight toward model-reading edges (B-type), which current column scores show still split the top agents.
- Rotate families across versions.

## Calibration
- Run a local panel of released agents of different owners: crown-v228, dan0l v59, beast-d41 v19, our ditto c228 tree.
- Use the real `ScoreCases` logic; about $0.3 per round of 3 packs.
- Target: most cases discriminating and few unsolved.
- Risk: the live panel is the current top agents (unreleased code), which may be past some code limits.
