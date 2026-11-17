"""Show cases of a pack with every record that writes one of its items or distractors, the asked name, or its first
name (to judge by eye whether a careful reader can answer).

    python3 -I tools/show_case.py PACK.json [N] [START]
"""
import json, sys
p = json.load(open(sys.argv[1]))
n = int(sys.argv[2]) if len(sys.argv) > 2 else 3
start = int(sys.argv[3]) if len(sys.argv) > 3 else 0
recs = [(w.get("wave", wi), x) for wi, w in enumerate(p["memory_waves"]) for x in w["pairs"]]
for c in p["cases"][start:start + n]:
    q = c["question"]
    print(f"\n== {c['question_type']} | {q}\n   items {c['answer_items']}  traps {c['distractor_answers']}")
    name = q.split(" does ")[-1].split(" go to")[0].split(" love")[0].split(" keeps")[0] if " does " in q else q.split(" the ")[-1]
    keys = c["answer_items"] + c["distractor_answers"]
    for wv, x in sorted(recs, key=lambda r: r[1].get("timestamp", "")):
        t = x["prompt"]
        if any(k in t for k in keys):
            tag = "EV" if x["pair_id"] in c["evidence_pair_ids"] else "  "
            print(f"   {tag} w{wv} {x.get('timestamp', '')[5:16]} {x['session_id']}: {t[:170]!r}")
