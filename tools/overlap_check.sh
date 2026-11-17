#!/bin/bash
# Sentences in our generator's families that also appear verbatim in a downloaded released generator's source
# (released sources under ditto/.cache/miner_gens/dl, vendor skipped). Prints each shared quoted string of
# 30+ characters and where it occurs.
D=/home/bittensor/ditto/.cache/miner_gens/dl
G=/home/bittensor/ditto-generator/generator/families
python3 - "$D" "$G" <<'P'
import glob, os, re, sys
D, G = sys.argv[1:]
q = re.compile(r'"([^"\\]{30,})"')
def strings(path):
    try:
        return set(q.findall(open(path, errors="ignore").read()))
    except Exception:
        return set()
ours = {}
for f in glob.glob(G + "/*.go"):
    for s in strings(f):
        ours.setdefault(s, os.path.basename(f))
norm = lambda s: re.sub(r"\{[^}]*\}|%[a-z]", "#", s.lower())
oursn = {norm(s): (s, f) for s, f in ours.items()}
hits = 0
for f in glob.glob(D + "/*/src/**/*.go", recursive=True):
    if "/vendor/" in f:
        continue
    for s in strings(f):
        k = norm(s)
        if k in oursn:
            hits += 1
            print(f"SAME  ours {oursn[k][1]}: {oursn[k][0]}\n      theirs {f.split('/dl/')[1]}")
print("shared strings:", hits)
P
grep -n "set foot\|lunch\|lift" $G/multi.go $G/edges.go | head
