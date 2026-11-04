"""Explain each panel case: why every agent passed or failed it, from the scorer's results and the traces.

    python3 tools/panel_cases.py practice/generator/<stamp>

For each case and agent run it puts the official score next to what the trace shows, and names the
cause of a miss:

    missing k/n    some evidence records never reached the model (retrieval)
    decoy X        the evidence was there but the reply named a distractor (reasoning or stale value)
    gold+decoy X   the reply named the right value and a distractor, which the grader scores 0
    declined       the reply was empty or said it could not tell
    other          the evidence was there and the reply named something else

It prints a cause count per family and writes cases.md into the run folder with every case the
agents disagree on (or all cases with --all), so a split can be traced to its record or its reply.
Reads local files only; costs nothing.
"""

import argparse
import collections
import json
import re
import sys
from pathlib import Path

DECLINE = re.compile(r"\b(i (?:do not|don't|cannot|can't) (?:know|tell|find|see)|not (?:sure|mentioned|specified)|"
                     r"no (?:record|information|mention)|unknown|unclear)\b", re.I)


def norm(text: str) -> str:
    return re.sub(r"[^a-z0-9@.]+", " ", (text or "").lower()).strip()


def mentions(reply: str, value: str) -> bool:
    items = [x for x in re.split(r",|\band\b", value) if x.strip()] or [value]
    padded = f" {norm(reply)} "
    return all(f" {norm(x)} " in padded for x in items)


def family(case_id: str) -> str:
    return re.sub(r"-\d+$", "", case_id)


def stem_rep(name: str, prefix: str) -> tuple[str, str]:
    """Split "run-<pack>-r0" (or "graded-…", or a trace folder "<pack>-r0") into pack stem and repeat."""
    stem, rep = name[len(prefix):].rsplit("-", 1)
    return stem, rep


def load(run_dir: Path):
    # Packs from different seeds reuse case IDs for different questions, so every case is keyed by
    # "<pack stem>/<case id>".
    cases, prompt, by_question = {}, {}, {}
    for pack_path in sorted((run_dir / "packs").glob("*.json")):
        stem, pack = pack_path.stem, json.loads(pack_path.read_text())
        wave_of = {}
        for w, wave in enumerate(pack["memory_waves"]):
            for pair in wave["pairs"]:
                prompt[f"{stem}/{pair['pair_id']}"] = pair["prompt"][:120]
                wave_of[pair["pair_id"]] = w
        for c in pack["cases"]:
            waves = sorted(wave_of[e] for e in c["evidence_pair_ids"] if e in wave_of)
            c["_stem"], c["_waves"] = stem, waves
            # A case run after wave 0 has no run_after_wave in the pack (the field is omitted when zero).
            c["_after"] = c.get("run_after_wave", 0)
            # How far back the oldest evidence sits from the wave the case runs after, and how many records.
            c["_span"] = f"back {c['_after'] - waves[0]}, {len(waves)} rec" if waves else "?"
            cases[f"{stem}/{c['id']}"] = c
        # A question asked again after a later wave is a different case; the scorer finishes a wave's
        # cases before seeding the next, so repeats appear in a trace in wave order.
        same = collections.defaultdict(list)
        for c in sorted(pack["cases"], key=lambda c: c["_after"]):
            same[norm(c["question"])].append(f"{stem}/{c['id']}")
        by_question[stem] = same

    scores = collections.defaultdict(dict)  # "<stem>/<case id>" -> {(agent, rep): score}
    status = {}
    for run in sorted(run_dir.glob("run-*.json")):
        stem, rep = stem_rep(run.stem, "run-")
        d = json.loads(run.read_text())
        for seed in d["input"]["seeds"]:
            for slot in seed["panel"]:
                for r in slot["results"]:
                    scores[f"{stem}/{r['case_id']}"][(slot["agent_id"], rep)] = float(r["score"])
        for seed in d["score"]["seeds"]:
            for c in seed["cases"]:
                status.setdefault(f"{stem}/{c['case_id']}", collections.Counter())[c["status"]] += 1
    # A run stopped before the scorer wrote its results: tools/gradetraces grades the saved replies.
    scored = {(cid.split("/")[0], rep) for cid, s in scores.items() for (_, rep) in s}
    for graded in sorted(run_dir.glob("graded-*.json")):
        stem, rep = stem_rep(graded.stem, "graded-")
        if (stem, rep) in scored:
            continue
        for agent, per_case in json.loads(graded.read_text())["scores"].items():
            for cid, score in per_case.items():
                scores[f"{stem}/{cid}"][(agent, rep)] = float(score)
    for cid, got in scores.items():
        if cid not in status:
            solved = sum(s >= 0.5 for s in got.values())
            status[cid] = collections.Counter({f"{solved}/{len(got)} runs solved": 1})

    traces = {}  # ("<stem>/<case id>", agent, rep) -> trace row
    for f in sorted((run_dir / "trace").glob("*/*/trace.jsonl")):
        (stem, rep), agent = stem_rep(f.parent.parent.name, ""), f.parent.name
        asked = collections.Counter()
        for line in f.read_text().splitlines():
            if line.strip():
                row = json.loads(line)
                q = norm(row.get("user_input"))
                same = by_question.get(stem, {}).get(q) or []
                if same:
                    traces[(same[min(asked[q], len(same) - 1)], agent, rep)] = row
                asked[q] += 1
    return cases, prompt, scores, status, traces


def explain(case: dict, prompt: dict, row: dict | None, score: float | None) -> tuple[str, str]:
    if row is None:
        return "no trace", ""
    reply = (row.get("final_text") or "").strip()
    seen = {p[1][:120] for p in row.get("packet") or []}
    ids = [f"{case['_stem']}/{e}" for e in case["evidence_pair_ids"]]
    evidence = [prompt[e] for e in ids if e in prompt]
    have = sum(e in seen for e in evidence)
    decoys = [d for d in case.get("distractor_answers") or [] if d and mentions(reply, d)]
    gold = mentions(reply, case["expected_answer"])
    if score is not None and score >= 0.5:
        cause = "pass"
    elif gold and decoys:
        cause = f"gold+decoy {decoys[0]}"
    elif not reply or (DECLINE.search(reply) and not gold):
        cause = "declined"
    elif have < len(evidence):
        cause = f"missing {len(evidence) - have}/{len(evidence)}"
    elif decoys:
        cause = f"decoy {decoys[0]}"
    elif score is None:
        cause = "unscored"
    else:
        cause = "other"
    return cause, f"{row.get('lane', '?')} {row.get('path', '?')} packet {have}/{len(evidence)}"


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("run_dir", type=Path)
    ap.add_argument("--all", action="store_true", help="write every case to cases.md, not only the split ones")
    args = ap.parse_args()
    if not (args.run_dir / "packs").is_dir():
        print(f"no packs/ in {args.run_dir}", file=sys.stderr)
        return 2
    cases, prompt, scores, status, traces = load(args.run_dir)
    runs = sorted({(a, r) for (_, a, r) in traces} | {k for s in scores.values() for k in s})
    agents = sorted({a for a, _ in runs})

    causes = collections.defaultdict(lambda: collections.defaultdict(collections.Counter))
    by_type = collections.defaultdict(lambda: collections.defaultdict(list))
    by_span = collections.defaultdict(lambda: collections.defaultdict(list))
    lines =[f"# Panel cases: {args.run_dir.name}", "",
             "Cause per agent run (`agent/rN`): pass, missing k/n (evidence not in the packet), "
             "decoy X, gold+decoy X, declined, other.", ""]
    for cid, case in cases.items():
        got = {}
        for agent, rep in runs:
            score = scores[cid].get((agent, rep))
            cause, where = explain(case, prompt, traces.get((cid, agent, rep)), score)
            got[(agent, rep)] = (score, cause, where)
            causes[family(case["id"])][agent][cause.split(" ")[0]] += 1
            if cause != "no trace":
                by_type[case["question_type"]][agent].append(cause == "pass")
                back = case.get("run_after_wave", 0) - case["_waves"][0] if case["_waves"] else -1
                far = "0-1" if back <= 1 else "2-3" if back <= 3 else "4-5" if back <= 5 else "6-7"
                many = len(case["_waves"])
                by_span[f"oldest evidence {far} waves back"][agent].append(cause == "pass")
                by_span[f"evidence in {many if many < 4 else '4+'} records"][agent].append(cause == "pass")
        passed = {a for (a, _), (s, c, _) in got.items() if c == "pass"}
        if not args.all and (not passed or passed == set(agents)):
            continue
        st = ", ".join(f"{k} {v}" for k, v in sorted(status.get(cid, {}).items())) or "not scored yet"
        lines += [f"## {cid}  ({st})", "",
                  f"- question: {case['question']}",
                  f"- gold: {case['expected_answer']}; decoys: {', '.join(case.get('distractor_answers') or [])}",
                  f"- type {case['question_type']}; evidence in waves {case['_waves']} ({case['_span']}), "
                  f"asked after wave {case['_after']}", ""]
        for (agent, rep), (score, cause, where) in got.items():
            reply = (traces.get((cid, agent, rep)) or {}).get("final_text") or ""
            shown = "-" if score is None else f"{score:.2f}"
            lines.append(f"- `{agent}/{rep}` {shown} **{cause}** ({where}): {reply.strip()[:160]!r}")
        lines.append("")

    print(f"{'family':34s} " + " ".join(f"{a[:22]:>24s}" for a in agents))
    for fam in sorted(causes):
        cells = []
        for a in agents:
            c = causes[fam][a]
            cells.append(" ".join(f"{k[:4]}{v}" for k, v in c.most_common()) or "-")
        print(f"{fam:34s} " + " ".join(f"{c:>24s}" for c in cells))
    for title, table in (("question type", by_type), ("evidence distance", by_span)):
        print(f"\npassed by {title}")
        for key in sorted(table):
            cells = [f"{sum(table[key][a])}/{len(table[key][a])}" if table[key][a] else "-" for a in agents]
            print(f"  {key:44s} " + " ".join(f"{c:>9s}" for c in cells))
    out = args.run_dir / "cases.md"
    out.write_text("\n".join(lines) + "\n")
    print(f"\nwrote {out}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
