"""Print the first differing case and memory pair between two packs.
    python3 -I tools/first_diff.py PACK_A PACK_B
"""
import json, sys
pa, pb = json.load(open(sys.argv[1])), json.load(open(sys.argv[2]))
print("top-level keys differing:", [k for k in pa if pa.get(k) != pb.get(k)])
for i, (x, y) in enumerate(zip(pa["cases"], pb["cases"])):
    if x != y:
        print("case", i, {k: (x[k], y.get(k)) for k in x if x[k] != y.get(k)})
        break
ra = [p for w in pa["memory_waves"] for p in w["pairs"]]
rb = [p for w in pb["memory_waves"] for p in w["pairs"]]
print("pairs", len(ra), len(rb))
for i, (x, y) in enumerate(zip(ra, rb)):
    if x != y:
        print("pair", i, {k: (str(x[k])[:160], str(y.get(k))[:160]) for k in x if x[k] != y.get(k)})
        break
