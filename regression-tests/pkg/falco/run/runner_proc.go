package run

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type procRunner struct {
	binaryPath string
}

// writeConfigToTempFile encodes a config to a newly-created temporary file
// and returns the file name or a non-nil error in case of failure.
// The newly-created file should be deleted manually from the filesystem
func (p *procRunner) writeConfigToTempFile(c *Config) (string, error) {
	f, err := ioutil.TempFile("", "falco-runner-")
	if err != nil {
		return "", err
	}
	name := f.Name()
	err = c.Marshal(f)
	if err == nil {
		err = f.Close()
		if err == nil {
			return name, nil
		}
	}
	os.Remove(name)
	return "", err
}

// NewProcRunner returns a Falco runner that runs a binary process
func NewProcRunner(binaryPath string) Runner {
	return &procRunner{binaryPath: binaryPath}
}

func (p *procRunner) Run(ctx context.Context, options ...RunnerOption) error {
	opts := &runnerOptions{
		config:  &Config{},
		options: []string{},
		stderr:  io.Discard,
		stdout:  io.Discard,
	}
	for _, o := range options {
		o(opts)
	}

	// create temp file to dump the YAML configuration
	confPath, err := p.writeConfigToTempFile(opts.config)
	if err != nil {
		return err
	}
	defer os.Remove(confPath)

	// launch Falco process
	opts.options = append(opts.options, "-c")
	opts.options = append(opts.options, confPath)
	cmd := exec.CommandContext(ctx, p.binaryPath, opts.options...)
	cmd.Stdout = opts.stdout
	cmd.Stderr = opts.stderr

	err = cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		err = &CodeError{Code: exitErr.ExitCode()}
	}
	return err
}
