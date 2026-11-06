// Command sandboxcheck runs a generator through qualification the way the live scorer does, in the
// scorer's Docker sandbox, before anything is uploaded.
//
//	go run . -dir ../../generator -cases 70 -bench 16
//
// It packs the directory exactly as `ditto generator upload` does (genrun.PackDir: the contract's file
// rules, dependencies vendored), unpacks the archive into a scratch directory opened up to 0755 as
// challenge-scorer does, and evaluates it with genrun.Evaluate in genrun.Docker: two scored seeds, a
// probe seed, validation of every pack and the determinism re-run, with no network, a read-only root,
// an unprivileged user and the scorer's limits. `challenge rehearse -sandbox docker` cannot do this
// when run as root: it leaves the scratch directory at MkdirTemp's 0700, which the container user
// cannot read, so every generator (the starter kit included) fails with "go.mod file not found".
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/genrun"
)

func main() {
	dir := flag.String("dir", "../../generator", "generator module directory")
	cases := flag.Int("cases", 70, "cases per seed")
	bench := flag.Int("bench", 16, "round bench version the requests are frozen at")
	first := flag.Int64("seed", 1001, "first scored seed (the probe is the next one after the scored seeds)")
	seeds := flag.Int("seeds", 2, "scored seeds, as in qualification")
	hotkey := flag.String("hotkey", "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY", "miner hotkey the packs echo")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	raw, info, err := genrun.PackDir(ctx, *dir)
	if err != nil {
		fail("pack: %v", err)
	}
	fmt.Printf("archive: %d bytes, vendored %v\n", len(raw), info.Vendored)

	src, err := os.MkdirTemp("", "sandboxcheck-src-*")
	if err != nil {
		fail("scratch dir: %v", err)
	}
	defer os.RemoveAll(src)
	if err := os.Chmod(src, 0o755); err != nil {
		fail("chmod: %v", err)
	}
	if _, err := challenge.ExtractArchive(raw, src); err != nil {
		fail("archive rejected by the contract: %v", err)
	}

	plan := genrun.Plan{MinerHotkey: *hotkey, Cases: *cases, CheckDeterminism: true, RoundBenchVersion: *bench}
	for i := int64(0); i < int64(*seeds); i++ {
		plan.Seeds = append(plan.Seeds, challenge.SeedFromInt64(*first+i))
	}
	plan.ProbeSeed = challenge.SeedFromInt64(*first + int64(*seeds))

	ev, err := genrun.Evaluate(ctx, genrun.Docker{}, src, plan)
	if err != nil {
		fail("sandbox (not the generator's fault): %v", err)
	}
	fmt.Printf("build: ok=%v %dms %s\n", ev.Build.OK, ev.Build.DurationMS, ev.Build.Error)
	for i, s := range ev.Seeds {
		report(fmt.Sprintf("seed %d", i+1), s)
	}
	if ev.Probe != nil {
		report("probe", *ev.Probe)
	}
	if ev.Deterministic != nil {
		fmt.Printf("deterministic: %v\n", *ev.Deterministic)
	}
	if ev.Valid() {
		fmt.Println("QUALIFIES: the scorer would accept this generator")
		return
	}
	fmt.Println("REJECTED: the scorer would refuse this generator")
	os.Exit(1)
}

func report(label string, s genrun.SeedOutput) {
	fmt.Printf("%s: valid=%v cases=%d policy=%d %dms", label, s.Valid(), len(s.Report.Cases), s.GradingBenchVersion, s.DurationMS)
	if s.Error != "" {
		fmt.Printf(" error=%s", s.Error)
	}
	if !s.Report.Valid && len(s.Report.Issues) > 0 {
		b, _ := json.Marshal(s.Report.Issues[:min(3, len(s.Report.Issues))])
		fmt.Printf(" issues=%s", b)
	}
	fmt.Println()
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}
