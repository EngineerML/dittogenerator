"""Group overlap_check.sh output: each of our shared strings once, with the released generators that hold it.
Skips import paths. Strings shared with the public starter kit are marked KIT.
    bash tools/overlap_check.sh > /tmp/o.txt; python3 tools/overlap_list.py /tmp/o.txt
"""
import collections, re, subprocess, sys
kit = open("/dev/null").read()
try:
    kit = subprocess.run("cat /home/bittensor/ditto-generator/repos/ditto-subnet/miners/dittobench-generator-starter-kit/families/*.go",
                         shell=True, capture_output=True, text=True).stdout
except Exception:
    pass
lines = open(sys.argv[1]).read().splitlines()
where = collections.defaultdict(set)
for a, b in zip(lines, lines[1:]):
    m = re.match(r"SAME  ours (\S+): (.*)", a)
    if m and b.strip().startswith("theirs") and "github.com" not in m.group(2):
        where[(m.group(1), m.group(2))].add(b.split("theirs ")[1].split("/")[0])
for (f, s), gens in sorted(where.items()):
    tag = "KIT " if s in kit else "    "
    print(f"{tag}{f:10} {s[:90]}  <- {', '.join(sorted(gens))}")
