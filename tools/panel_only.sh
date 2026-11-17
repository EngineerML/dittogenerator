#!/bin/bash
# Paid 4-agent panel on already-built packs (.cache/packs/TAG), then the trap table. For a retry after a
# round_v81.sh whose panel step failed.
#   bash tools/panel_only.sh TAG
T=${1:?tag}
bash /home/bittensor/ditto/practice/gen_v8_panel.sh "$T"
cd /home/bittensor/ditto && python3 practice/gen_trap_table.py "$(ls -td practice/generator/2026*Z | head -1)"
