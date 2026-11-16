#!/bin/bash
# Build, vet and test the generator, then rehearse it as a bench-18 round (3 seeds, 35 and 80 cases) with the
# competition's own tool.
export GOTOOLCHAIN=auto GOFLAGS=-mod=mod
cd /home/bittensor/ditto-generator/generator || exit 1
go vet ./... 2>&1 | tail -10 && echo "vet ok"
go test ./... 2>&1 | tail -15
cd /home/bittensor/ditto-generator/repos/ditto-subnet/research/dittobench-datagen || exit 1
for n in 35 80; do
  go run ./cmd/challenge rehearse -dir ../../../../generator -hotkey 5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY -seeds 3 -cases $n -round-bench 18 2>&1 | tail -12
done
