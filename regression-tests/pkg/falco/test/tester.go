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
	config   string
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

// Condition is an assertion over a given Falco run
type Condition func(*testerState) error

func setupTestBench(tb testing.TB, runner run.Runner, conditions ...Condition) (*testerState, []run.RunnerOption) {
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

	return state, []run.RunnerOption{
		run.WithStdout(&state.stdout),
		run.WithStderr(&state.stderr),
		run.WithConfig(state.config),
		run.WithArgs(state.args),
	}
}

// RunTest runs a test with the given conditions with the given runner
func RunTest(t *testing.T, runner run.Runner, conditions ...Condition) {
	state, runOpts := setupTestBench(t, runner, conditions...)
	// todo: add some drift to the duration
	ctx, cancel := context.WithTimeout(context.Background(), state.duration)
	defer cancel()
	state.err = runner.Run(ctx, runOpts...)
	state.done = true

	// run post-run conditions
	for _, c := range conditions {
		if err := c(state); err != nil {
			t.Fatal(err.Error())
		}
	}
}

// RunBench runs a benchmark with the given conditions with the given runner
func RunBench(b *testing.B, runner run.Runner, conditions ...Condition) {
	state, runOpts := setupTestBench(b, runner, conditions...)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		state.stderr.Reset()
		state.stdout.Reset()
		// todo: should we catch errors?
		runner.Run(ctx, runOpts...)
	}
}
