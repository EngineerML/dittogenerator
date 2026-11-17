"""Print every record (wave, time, prompt) of a pack that contains a phrase, case-insensitively.
    python3 -I tools/grep_pack.py PACK.json PHRASE
"""
import json, sys
p = json.load(open(sys.argv[1]))
want = sys.argv[2].lower()
for wi, w in enumerate(p["memory_waves"]):
    for x in w["pairs"]:
        if want in x["prompt"].lower():
            print(wi, x.get("id") or x.get("pair_id"), repr(x["prompt"][:300]))
