#!/bin/bash
# Fast-forward repos/ditto-subnet from the synced /home/bittensor/ditto/repos copy (never from GitHub), then rebuild
# .cache/challenge-scorer from services/dittobench-api, and build the generator once to check it still compiles.
source /etc/profile.d/zz-proxy.sh
cd /home/bittensor/ditto-generator/repos/ditto-subnet || exit 1
git fetch -q synced main && git -c remote.origin.promisor=false merge -q --ff-only synced/main && git update-ref refs/remotes/origin/main synced/main
git log --oneline -1
cd services/dittobench-api || exit 1
export GOTOOLCHAIN=auto
go build -o /home/bittensor/ditto-generator/.cache/challenge-scorer.new ./cmd/challenge-scorer 2>&1 | tail -5 && mv /home/bittensor/ditto-generator/.cache/challenge-scorer.new /home/bittensor/ditto-generator/.cache/challenge-scorer && echo "scorer rebuilt"
cd /home/bittensor/ditto-generator/generator && go vet ./... 2>&1 | tail -5 && go build -o /tmp/horizon-gen . && echo "generator builds"
