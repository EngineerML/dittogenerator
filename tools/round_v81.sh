#!/bin/bash
# Calibration round: rehearse (35, 80), sandbox packs into .cache/packs/TAG, then the paid 4-agent panel on seed-1
# and the trap table.
#   bash tools/round_v81.sh TAG
T=${1:-v81}
bash /home/bittensor/ditto-generator/tools/rehearse18.sh 35 80 2>&1 | grep -E 'valid|invalid|vet'
bash /home/bittensor/ditto-generator/tools/packs_v8.sh $T 2>&1 | grep -E 'QUALIF|REJECT|seed [0-9]:|Counter'
bash /home/bittensor/ditto/practice/gen_v8_panel.sh $T
cd /home/bittensor/ditto && python3 practice/gen_trap_table.py $(ls -td practice/generator/2026*Z | head -1)
