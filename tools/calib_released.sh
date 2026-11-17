#!/bin/bash
# Calibrate the local panel on a released generator: build its packs in the scorer's Docker sandbox (bench 18,
# 35 cases, 3 seeds) into .cache/packs/TAG, run the paid 4-agent panel on seed-1, then per-family case status
# from the scorer's own result.
#   bash tools/calib_released.sh SRC_DIR TAG
SRC=${1:?source dir}; T=${2:?tag}
export GOTOOLCHAIN=auto
out=/home/bittensor/ditto-generator/.cache/packs/$T
rm -rf "$out"; mkdir -p "$out"
cd /home/bittensor/ditto-generator/tools/sandboxcheck || exit 1
go run . -dir "$SRC" -cases 35 -bench 18 -seeds 3 -out-dir "$out" 2>&1 | tail -6
[ -f "$out/seed-1.json" ] || { echo "no seed-1 pack, stopping before the paid panel"; exit 1; }
bash /home/bittensor/ditto/practice/gen_v8_panel.sh "$T"
cd /home/bittensor/ditto && python3 practice/gen_case_status.py "$(ls -td practice/generator/2026*Z | head -1)"
