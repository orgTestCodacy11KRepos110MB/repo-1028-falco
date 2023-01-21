package run

import (
	"context"
	"fmt"
	"io"
)

// RunnerOption is an option for running Falco
type RunnerOption func(*runnerOptions)

// Runner runs Falco with a given set of options
type Runner interface {
	// Run runs Falco with the given options and returns when it finishes its
	// execution or when the context deadline is exceeded.
	// Returns a non-nil error in case of failure.
	Run(ctx context.Context, options ...RunnerOption) error
}

type runnerOptions struct {
	config  string
	options []string
	stderr  io.Writer
	stdout  io.Writer
}

func WithConfig(config string) RunnerOption {
	return func(ro *runnerOptions) { ro.config = config }
}

func WithOptions(options []string) RunnerOption {
	return func(ro *runnerOptions) { ro.options = options }
}

func WithStdout(writer io.Writer) RunnerOption {
	return func(ro *runnerOptions) { ro.stdout = writer }
}

func WithStderr(writer io.Writer) RunnerOption {
	return func(ro *runnerOptions) { ro.stderr = writer }
}

// CodeError is an error represented by a numeric code
type CodeError struct {
	Code int
}

func (c *CodeError) Error() string {
	return fmt.Sprintf("error code %d", c.Code)
}
