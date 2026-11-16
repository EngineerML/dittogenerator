#!/bin/bash
# Build v8 packs in the scorer's Docker sandbox (bench 18, 35 cases, 3 seeds) into .cache/packs/v8, then show the
# family mix of the first pack.
export GOTOOLCHAIN=auto
cd /home/bittensor/ditto-generator/tools/sandboxcheck || exit 1
out=/home/bittensor/ditto-generator/.cache/packs/${1:-v8}
rm -rf $out; mkdir -p $out
go run . -dir ../../generator -cases 35 -bench 18 -seeds 3 -out-dir $out 2>&1 | tail -8
ls -la $out
python3 - $out <<'P'
import json, sys, glob, collections
for f in sorted(glob.glob(sys.argv[1] + "/*.json"))[:1]:
    p = json.load(open(f))
    print(f, "grading", p["grading_bench_version"], "pairs", sum(len(w["pairs"]) for w in p["memory_waves"]))
    print(collections.Counter(c["question_type"] for c in p["cases"]))
    for c in p["cases"][:6]:
        print(" ", c["question_type"], "|", c["question"], "|", c["answer_items"], "| traps", c["distractor_answers"])
P
