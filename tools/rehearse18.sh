#!/bin/bash
# Rehearse the generator as a bench-18 round (3 seeds) at the given case counts; no go test (that takes minutes).
#   bash tools/rehearse18.sh [35 80]
export GOTOOLCHAIN=auto GOFLAGS=-mod=mod
cd /home/bittensor/ditto-generator/generator && go vet ./... && echo "vet ok" || exit 1
cd /home/bittensor/ditto-generator/repos/ditto-subnet/research/dittobench-datagen || exit 1
for n in ${@:-35 80}; do
  go run ./cmd/challenge rehearse -dir ../../../../generator -hotkey 5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY -seeds 3 -cases $n -round-bench 18 2>&1 | tail -6
done
