package test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
)

const (
	// DefaultDuration is the default max duration of a Falco run
	DefaultDuration = time.Second * 5
)

type testerState struct {
	duration time.Duration
	done     bool
	err      error
	args     []string
	config   run.Config
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

// Condition is an assertion over a given Falco run
type Condition func(*testerState) error

// internal implementation for tests and benchmarks
func testBenchRun(tb testing.TB, runner run.Runner, conditions ...Condition) {
	// todo: set default configuration fields
	// todo: set default forced options (depends on the condition I guess)
	state := &testerState{
		duration: DefaultDuration,
	}

	// run pre-run conditions
	for _, c := range conditions {
		if err := c(state); err != nil {
			tb.Fatal(err.Error())
		}
	}

	// run Falco
	var runOpts []run.RunnerOption
	runOpts = append(runOpts, run.WithStdout(&state.stdout))
	runOpts = append(runOpts, run.WithStderr(&state.stderr))
	runOpts = append(runOpts, run.WithConfig(&state.config))
	runOpts = append(runOpts, run.WithOptions(state.args))
	ctx, cancel := context.WithTimeout(context.Background(), state.duration) // todo: add some drift margin
	defer cancel()
	state.err = runner.Run(ctx, runOpts...)
	state.done = true

	// run post-run conditions
	for _, c := range conditions {
		if err := c(state); err != nil {
			tb.Fatal(err.Error())
		}
	}
}

// TestProcess runs a test with the given conditions with the given runner
func TestProcess(t *testing.T, runner run.Runner, conditions ...Condition) {
	testBenchRun(t, runner, conditions...)
}

// BenchProcess runs a benchmark with the given conditions with the given runner
func BenchProcess(b *testing.B, runner run.Runner, conditions ...Condition) {
	testBenchRun(b, runner, conditions...)
}
