package run

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type procRunner struct {
	executable string
}

// writeToTempFile encodes a config to a newly-created temporary file
// and returns the file name and a callback for deleting the file,
// or a non-nil error in case of failure. The newly-created file should be
// deleted manually by invoking the returned callback.
func (p *procRunner) writeToTempFile(c string) (string, func() error, error) {
	f, err := os.CreateTemp("", "falco-runner-")
	if err != nil {
		return "", nil, err
	}
	name := f.Name()
	n, err := f.WriteString(c)
	if err == nil || n < len(c) {
		err = f.Close()
		if err == nil {
			return name, func() error { return os.Remove(name) }, nil
		}
		if n < len(c) {
			err = io.ErrShortWrite
		}
	}
	return "", nil, err
}

// NewExecutableRunner returns a Falco runner that runs a binary process
func NewExecutableRunner(executable string) Runner {
	return &procRunner{executable: executable}
}

func (p *procRunner) Run(ctx context.Context, options ...RunnerOption) error {
	opts := &runnerOptions{
		config: "",
		args:   []string{},
		stderr: io.Discard,
		stdout: io.Discard,
	}
	for _, o := range options {
		o(opts)
	}

	// create temp file to dump the YAML configuration
	conf, confDelete, err := p.writeToTempFile(opts.config)
	if err != nil {
		return err
	}
	defer confDelete()

	// launch Falco process
	opts.args = append(opts.args, "-c")
	opts.args = append(opts.args, conf)
	cmd := exec.CommandContext(ctx, p.executable, opts.args...)
	cmd.Stdout = opts.stdout
	cmd.Stderr = opts.stderr

	err = cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		err = &CodeError{Code: exitErr.ExitCode()}
	}
	return err
}
