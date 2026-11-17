#!/bin/bash
# Free checks on generator/: vet, unit tests, then the two-link plan sample (first cases of seed-like builds).
export GOTOOLCHAIN=auto
cd /home/bittensor/ditto-generator/generator || exit 1
go vet ./... && go test ./... 2>&1 | tail -5
