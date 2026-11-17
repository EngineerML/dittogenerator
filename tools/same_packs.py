"""Compare two pack directories built by packs_v8.sh: cases and memory waves equal per seed (generator version and
other provenance ignored).
    python3 -I tools/same_packs.py DIR_A DIR_B
"""
import json, sys, glob, os
a, b = sys.argv[1], sys.argv[2]
for f in sorted(glob.glob(a + "/*.json")):
    pa, pb = json.load(open(f)), json.load(open(os.path.join(b, os.path.basename(f))))
    same = pa["cases"] == pb["cases"] and pa["memory_waves"] == pb["memory_waves"]
    print(os.path.basename(f), "same" if same else "DIFFERENT")
