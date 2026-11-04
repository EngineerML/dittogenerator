// Command gradetraces grades the answers saved in a panel run's traces with the official memory grader.
//
//	go run . ../../practice/generator/<stamp>
//
// The challenge scorer only writes its results once every agent has answered, so a run stopped part way
// leaves paid answers with no scores. Each trace row holds the reply the harness returned; this grades it
// with grade.Memory under the policy the pack declares, which is the call the scorer makes, and writes
// graded-<pack>-rN.json (agent -> case id -> score) next to the run for tools/panel_cases.py.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ditto-assistant/dittobench-datagen/challenge"
	"github.com/ditto-assistant/dittobench-datagen/grade"
	"github.com/ditto-assistant/dittobench-datagen/protocol"
)

var nonWord = regexp.MustCompile(`[^a-z0-9@.]+`)

func norm(s string) string {
	return strings.TrimSpace(nonWord.ReplaceAllString(strings.ToLower(s), " "))
}

type traceRow struct {
	UserInput string `json:"user_input"`
	FinalText string `json:"final_text"`
	Answer    any    `json:"answer"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: gradetraces <run dir>")
		os.Exit(2)
	}
	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string) error {
	packs, _ := filepath.Glob(filepath.Join(dir, "packs", "*.json"))
	for _, path := range packs {
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pack, err := challenge.Decode(raw)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		// A question asked again after a later wave is a different case; the scorer finishes a wave's
		// cases before seeding the next, so repeats appear in a trace in wave order.
		index := map[string][]int{}
		order := make([]int, len(pack.Cases))
		for i := range order {
			order[i] = i
		}
		sort.SliceStable(order, func(x, y int) bool { return pack.Cases[order[x]].RunAfterWave < pack.Cases[order[y]].RunAfterWave })
		for _, i := range order {
			q := norm(pack.Cases[i].Question)
			index[q] = append(index[q], i)
		}
		stem := strings.TrimSuffix(filepath.Base(path), ".json")
		reps, _ := filepath.Glob(filepath.Join(dir, "trace", stem+"-r*"))
		for _, rep := range reps {
			scores := map[string]map[string]float64{}
			traces, _ := filepath.Glob(filepath.Join(rep, "*", "trace.jsonl"))
			for _, trace := range traces {
				agent := filepath.Base(filepath.Dir(trace))
				got, err := gradeTrace(pack, index, trace)
				if err != nil {
					return err
				}
				if len(got) > 0 {
					scores[agent] = got
				}
			}
			out := filepath.Join(dir, "graded-"+filepath.Base(rep)+".json")
			raw, _ := json.MarshalIndent(map[string]any{"grading_bench_version": pack.GradingBenchVersion, "scores": scores}, "", " ")
			if err := os.WriteFile(out, raw, 0o644); err != nil {
				return err
			}
			agents := make([]string, 0, len(scores))
			for a := range scores {
				agents = append(agents, a)
			}
			sort.Strings(agents)
			for _, a := range agents {
				solved, sum := 0, 0.0
				for _, s := range scores[a] {
					sum += s
					if s >= 0.5 {
						solved++
					}
				}
				fmt.Printf("%s %-8s %d/%d solved, mean %.3f (policy %d)\n", filepath.Base(rep), a, solved, len(scores[a]),
					sum/float64(len(scores[a])), pack.GradingBenchVersion)
			}
		}
	}
	return nil
}

func gradeTrace(pack challenge.Pack, index map[string][]int, path string) (map[string]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scores := map[string]float64{}
	asked := map[string]int{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var row traceRow
		if json.Unmarshal(sc.Bytes(), &row) != nil {
			continue
		}
		q := norm(row.UserInput)
		same := index[q]
		if len(same) == 0 {
			continue
		}
		i := same[min(asked[q], len(same)-1)]
		asked[q]++
		answer, _ := row.Answer.(string)
		v := grade.Memory(pack.GradingCase(i), protocol.RunResponse{FinalText: row.FinalText, Answer: answer})
		scores[pack.Cases[i].ID] = v.Score
	}
	return scores, sc.Err()
}
