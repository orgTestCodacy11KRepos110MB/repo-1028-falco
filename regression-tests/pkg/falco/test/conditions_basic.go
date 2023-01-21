package test

import (
	"fmt"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
)

func Config(c string) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			ts.config = c
		}
		return nil
	}
}

func Args(args ...string) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			ts.args = append(ts.args, args...)
		}
		return nil
	}
}

func ExitCode(f func(int) error) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			return nil
		}
		if ts.err == nil {
			return f(0)
		}
		if codeErr, ok := ts.err.(*run.CodeError); ok {
			return f(codeErr.Code)
		}
		return fmt.Errorf("could not retrieve exit code")
	}
}

func Stdout(f func(string) error) Condition {
	return func(ts *testerState) error {
		if ts.done {
			return f(ts.stdout.String())
		}
		return nil
	}
}

func Stderr(f func(string) error) Condition {
	return func(ts *testerState) error {
		if ts.done {
			return f(ts.stderr.String())
		}
		return nil
	}
}
