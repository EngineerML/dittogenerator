"""Estimate a live round score for a family mix from local panel runs, with the panel's strongest agent modelled.

    python3 tools/panel_estimate.py RUN_DIR [RUN_DIR ...] [--solver 0.85] [--champion hades26] [--cps 70]

The local panel lacks the live panel's strongest agent (dan0l v42 solved 0.89 of the top generators' shared
cases in round 26078; the other four sat at 0.28-0.46), so a local run judges "solvable" against too few
owners. This reads every agent's official scores per family from the run folders (run-*.json, or
graded-*.json from tools/gradetraces), adds a modelled solver that reaches the threshold with probability
--solver (on a list it earns 0.5 or 1, elsewhere 1), and scores sampled packs with the challenge scorer's
rules: at least two owners at or above 0.5, trivial when every agent is, credit 1 - mean times 1.5 when the
champion misses, and the squared unsolved penalty. The mix is the family counts of the newest pack in the
newest run folder, scaled to --cps cases, or --mix FAMILY=N,... Reads local files only; costs nothing.
"""

import argparse
import collections
import json
import random
import re
import sys
from pathlib import Path

PATHS = ("records", "program", "topic")


def family(case_id: str) -> str:
    """Case id without its round number and reading path: 'typo-contact-typo-topic-3' -> 'typo-contact-typo'."""
    base = re.sub(r"-\d+$", "", case_id)
    parts = base.rsplit("-", 1)
    return parts[0] if len(parts) == 2 and parts[1] in PATHS else base


def load_scores(run_dirs: list[Path]) -> dict[str, dict[str, list[float]]]:
    """family -> agent -> every official score seen for it."""
    out = collections.defaultdict(lambda: collections.defaultdict(list))
    for run_dir in run_dirs:
        scored = set()
        for run in sorted(run_dir.glob("run-*.json")):
            scored.add(run.stem[len("run-"):])
            for seed in json.loads(run.read_text())["input"]["seeds"]:
                for slot in seed["panel"]:
                    for r in slot["results"]:
                        out[family(r["case_id"])][slot["agent_id"]].append(float(r["score"]))
        for graded in sorted(run_dir.glob("graded-*.json")):
            if graded.stem[len("graded-"):] in scored:
                continue
            for agent, per_case in json.loads(graded.read_text())["scores"].items():
                for cid, score in per_case.items():
                    out[family(cid)][agent].append(float(score))
    return out


def judge(scores: list[float], champion: int | None) -> tuple[str, float]:
    """One case's verdict and credit under the default ScoreParams, every agent its own owner."""
    solvers = sum(s >= 0.5 for s in scores)
    if solvers < 2:
        return "unsolved", 0.0
    if solvers == len(scores):
        return "trivial", 0.0
    credit = 1 - sum(scores) / len(scores)
    if champion is not None and scores[champion] < 0.5:
        credit *= 1.5
    return "discriminating", credit


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("run_dirs", type=Path, nargs="+")
    ap.add_argument("--agents", default="hades26,horizon9,drp", help="local agents to use, comma-separated")
    ap.add_argument("--solver", type=float, default=0.85, help="chance the modelled strongest agent reaches 0.5")
    ap.add_argument("--champion", default="solver",
                    help="agent holding the crown: a local agent, 'solver' for the modelled one, '' for none")
    ap.add_argument("--cps", type=int, default=70, help="cases per seed in the live round")
    ap.add_argument("--mix", help="FAMILY=N,... instead of the newest pack's counts")
    ap.add_argument("--trials", type=int, default=4000)
    args = ap.parse_args()

    agents = [a for a in args.agents.split(",") if a]
    data = load_scores(args.run_dirs)
    if args.mix:
        mix = {k: int(v) for k, v in (x.split("=") for x in args.mix.split(","))}
    else:
        packs = sorted(args.run_dirs[-1].glob("packs/*.json"))
        if not packs:
            print("no pack to take the mix from; pass --mix", file=sys.stderr)
            return 2
        mix = collections.Counter(family(c["id"]) for c in json.loads(packs[-1].read_text())["cases"])
    total = sum(mix.values())
    # "solver" makes the modelled strongest agent the champion (dan0l v42 holds the crown since 2026-11-05).
    champion = len(agents) if args.champion == "solver" else agents.index(args.champion) if args.champion in agents else None
    rng = random.Random(7)
    lists = {"alias-list", "membership-fold", "membership-fold-readd"}

    def draw(fam: str) -> list[float] | None:
        per = data.get(fam)
        if not per or any(not per.get(a) for a in agents):
            return None
        scores = [rng.choice(per[a]) for a in agents]
        hit = rng.random() < args.solver
        scores.append((rng.choice([0.5, 1.0]) if fam in lists else 1.0) if hit else 0.0)
        return scores

    print(f"panel: {', '.join(agents)} + modelled solver at {args.solver:.0%}; champion {args.champion or 'none'}")
    print(f"{'family':28s} {'cases':>5s} " + " ".join(f"{a[:9]:>9s}" for a in agents) + "   disc  triv unsol  credit/case")
    missing, per_family = [], {}
    for fam, n in sorted(mix.items(), key=lambda kv: -kv[1]):
        if not data.get(fam) or any(not data[fam].get(a) for a in agents):
            missing.append(fam)
            continue
        verdicts, credit = collections.Counter(), 0.0
        for _ in range(args.trials):
            v, c = judge(draw(fam), champion)
            verdicts[v] += 1
            credit += c
        per_family[fam] = (verdicts, credit / args.trials)
        rates = " ".join(f"{sum(s >= 0.5 for s in data[fam][a]) / len(data[fam][a]):9.0%}" for a in agents)
        t = args.trials
        print(f"{fam:28s} {n:5d} {rates}   {verdicts['discriminating'] / t:4.0%} {verdicts['trivial'] / t:5.0%} "
              f"{verdicts['unsolved'] / t:5.0%}  {credit / t:6.2f}")
    if missing:
        print("no local measurement (left out of the estimate):", ", ".join(missing))

    # Whole-pack score: sample cps cases from the measured part of the mix, score them as one pack.
    fams = [f for f in mix if f in per_family]
    weights = [mix[f] for f in fams]
    scores = []
    for _ in range(400):
        credit, unsolved = 0.0, 0
        for fam in rng.choices(fams, weights, k=args.cps):
            v, c = judge(draw(fam), champion)
            credit += c
            unsolved += v == "unsolved"
        solvable = 1 - unsolved / args.cps
        scores.append(credit / (1.5 * args.cps) * solvable * solvable)
    scores.sort()
    measured = sum(weights)
    print(f"\nestimated live pack score: {sum(scores) / len(scores):.3f} "
          f"(10-90%: {scores[len(scores) // 10]:.3f}-{scores[9 * len(scores) // 10]:.3f}); "
          f"{measured} of {total} cases in the mix measured")
    return 0


if __name__ == "__main__":
    sys.exit(main())
